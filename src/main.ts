import * as core from '@actions/core'
import * as github from '@actions/github'
import { GitHub } from '@actions/github/lib/utils'

async function run(): Promise<void> {
  try {
    core.debug(JSON.stringify(github.context))

    const token = core.getInput('token', { required: true })
    const quiet = core.getInput('quiet', { required: false }) === 'true'
    const titleRegexInput = core.getInput('titleRegex', { required: false }) || `^\[([A-Z]{2,})(-)(\d+)\]` //eslint-disable-line no-useless-escape
    const bodyRegexInput = core.getInput('bodyRegex', { required: false }) || `\[([A-Z]{2,})(-)(\d+)\]` //eslint-disable-line no-useless-escape
    const noTicketInput = core.getInput('noTicket', { required: false }) || '[no-ticket]'

    const client = github.getOctokit(token)
    const prTitle: string = github.context?.payload?.pull_request?.title
    const prBody = github.context?.payload?.pull_request?.body

    if (prBody === undefined) {
      core.debug('pr body is undefined')
      core.setFailed('Could not retrieve the Pull Request body')
      return
    }

    const reTitle = new RegExp(titleRegexInput, 'g')
    const reBody = new RegExp(bodyRegexInput, 'g')
    const foundTixInTitle = reTitle.test(prTitle)
    const foundTixInBody = reBody.test(prBody)
    const foundNoTixInTitle = prTitle.startsWith(noTicketInput)
    const foundNoTixInBody = prBody.includes(noTicketInput)
    core.debug(`found tix in title: ${foundTixInTitle}`)
    core.debug(`found tix in body : ${foundTixInBody}`)
    core.debug(`found no-tix in title: ${foundNoTixInTitle}`)
    core.debug(`found no-tix in body : ${foundNoTixInBody}`)

    if (foundTixInTitle && foundTixInBody) {
      core.debug('success, found tix in both title and body')
      return
    }
    if (foundNoTixInTitle && foundNoTixInBody) {
      core.debug(`success, found ${noTicketInput} in both title and body`)
      return
    }

    core.setFailed('missing JIRA ticket in both PR title AND body')
    makePrComment(client, quiet)
  } catch (error) {
    core.setFailed(error.message)
  }
}

async function makePrComment(client: InstanceType<typeof GitHub>, quiet: boolean): Promise<void> {
  try {
    if (quiet) {
      core.debug(`quiet mode enabled. skipping creating PR comment.`)
      return
    }

    const { owner, repo, number } = github.context.issue
    await client.rest.pulls.createReview({
      owner,
      repo,
      pull_number: number,
      body: 'PR title AND body does not contain a reference to a JIRA ticket.',
      event: 'COMMENT'
    })
  } catch (error) {
    core.setFailed(`Failed to update PR`)
  }
}

run()
