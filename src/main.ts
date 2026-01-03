import * as core from '@actions/core'
import { context, getOctokit } from '@actions/github'
import { env } from 'process'

export const defaultTitleBodyRegex = /\[([a-zA-Z]{2,}-\d+)\]/
const defaultNoTicket = '[no-ticket]'

export async function run(): Promise<void> {
  try {
    core.info('starting pr lint')
    // eslint-disable-next-line no-console
    console.info('starting pr lint')
    core.debug(JSON.stringify(context))

    const isCI = env.IS_CI === 'true'
    const token = core.getInput('token', { required: true })
    const makePrComment = core.getInput('make_pr_comment', { required: false }) === 'true'
    const titleRegexInput = core.getInput('title_regex', { required: false }) || defaultTitleBodyRegex
    const bodyRegexInput = core.getInput('body_regex', { required: false }) || defaultTitleBodyRegex
    const noTicketInput = core.getInput('no_ticket', { required: false }) || defaultNoTicket
    core.info(`input make_pr_comment: ${makePrComment}`)
    core.info(`input title_regex    : ${titleRegexInput}`)
    core.info(`input body_regex     : ${bodyRegexInput}`)
    core.info(`input no_ticket      : ${noTicketInput}`)

    const client = getOctokit(token)
    const prTitle: string = context.payload.pull_request?.title || ''
    const prBody = context.payload.pull_request?.body || ''

    const reTitle = new RegExp(titleRegexInput)
    const reBody = new RegExp(bodyRegexInput)
    const foundTixInTitle = reTitle.test(prTitle)
    const foundTixInBody = reBody.test(prBody)
    const foundNoTixInTitle = prTitle.includes(noTicketInput)
    const foundNoTixInBody = prBody.includes(noTicketInput)
    core.info(`found ticket in title   : ${foundTixInTitle}`)
    core.info(`found ticket in body    : ${foundTixInBody}`)
    core.info(`found no_ticket in title: ${foundNoTixInTitle}`)
    core.info(`found no_ticket in body : ${foundNoTixInBody}`)

    if (context.eventName !== 'pull_request') {
      core.info('success, event is not a pull request')
      return
    }
    if (foundTixInTitle && foundTixInBody) {
      core.info('success, found tix in both title and body')
      return
    }
    if (foundNoTixInTitle && foundNoTixInBody) {
      core.info(`success, found ${noTicketInput} in both title and body`)
      return
    }
    if (isCI) {
      core.info('ci mode detected. not returning failure')
      return
    }

    core.setFailed('missing ticket in both PR title AND body')
    if (makePrComment === true) {
      createPRComment(client, context)
    }
  } catch (e) {
    if (e instanceof Error) core.setFailed(e.message)
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
    core.error(`Failed to create PR comment: ${e}`)
  }
}
