import { debug, info, error, getInput, setFailed } from '@actions/core'
import { context, getOctokit } from '@actions/github'
import { env } from 'process'

export const defaultTitleBodyRegex = /\[([a-zA-Z]{2,}-\d+)\]/
const defaultNoTicket = '[no-ticket]'

export async function run(): Promise<void> {
  try {
    debug(JSON.stringify(context))

    const isCI = env.IS_CI === 'true'
    const token = getInput('token', { required: true })
    const makePrComment = getInput('make_pr_comment', { required: false }) === 'true'
    const titleRegexInput = getInput('title_regex', { required: false }) || defaultTitleBodyRegex
    const bodyRegexInput = getInput('body_regex', { required: false }) || defaultTitleBodyRegex
    const noTicketInput = getInput('no_ticket', { required: false }) || defaultNoTicket
    info(`input make_pr_comment: ${makePrComment}`)
    info(`input title_regex    : ${titleRegexInput}`)
    info(`input body_regex     : ${bodyRegexInput}`)
    info(`input no_ticket      : ${noTicketInput}`)

    const client = getOctokit(token)
    const prTitle: string = context.payload.pull_request?.title || ''
    const prBody = context.payload.pull_request?.body || ''

    const reTitle = new RegExp(titleRegexInput)
    const reBody = new RegExp(bodyRegexInput)
    const foundTixInTitle = reTitle.test(prTitle)
    const foundTixInBody = reBody.test(prBody)
    const foundNoTixInTitle = prTitle.includes(noTicketInput)
    const foundNoTixInBody = prBody.includes(noTicketInput)
    info(`found ticket in title   : ${foundTixInTitle}`)
    info(`found ticket in body    : ${foundTixInBody}`)
    info(`found no_ticket in title: ${foundNoTixInTitle}`)
    info(`found no_ticket in body : ${foundNoTixInBody}`)

    if (context.eventName !== 'pull_request') {
      info('success, event is not a pull request')
      return
    }
    if (foundTixInTitle && foundTixInBody) {
      info('success, found tix in both title and body')
      return
    }
    if (foundNoTixInTitle && foundNoTixInBody) {
      info(`success, found ${noTicketInput} in both title and body`)
      return
    }
    if (isCI) {
      info('ci mode detected. not returning failure')
      return
    }

    setFailed('missing ticket in both PR title AND body')
    if (makePrComment === true) {
      createPRComment(client, context)
    }
  } catch (e) {
    if (e instanceof Error) setFailed(e.message)
  }
}

async function createPRComment(client: ReturnType<typeof getOctokit>, ctx: typeof context): Promise<void> {
  try {
    // https://octokit.github.io/rest.js/v18#pulls-create-review-comment
    await client.rest.issues.createComment({
      owner: ctx.issue.owner,
      repo: ctx.issue.repo,
      issue_number: ctx.issue.number,
      body: 'PR title AND body does not contain a reference to a ticket.',
    })
  } catch (e) {
    error(`Failed to create PR comment: ${e}`)
  }
}
