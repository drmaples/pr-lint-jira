/* eslint no-console: */

import * as core from '@actions/core'
import * as github from '@actions/github'
import * as process from 'process'
import { GitHub } from '@actions/github/lib/utils'

async function run(): Promise<void> {
  try {
    core.debug(JSON.stringify(github.context))

    const isCI = process.env.IS_CI === 'true'
    const token = core.getInput('token', { required: true })
    const makePrComment = core.getInput('makePrComment', { required: false }) === 'true'
    const titleRegexInput = core.getInput('titleRegex', { required: false }) || `^\\[([A-Z]{2,}-\\d{3,})\\]`
    const bodyRegexInput = core.getInput('bodyRegex', { required: false }) || `\\[([A-Z]{2,}-\\d{3,})\\]`
    const noTicketInput = core.getInput('noTicket', { required: false }) || '[no-ticket]'

    const client = github.getOctokit(token)
    const prTitle: string = github.context?.payload?.pull_request?.title || ''
    const prBody = github.context?.payload?.pull_request?.body || ''

    const reTitle = new RegExp(titleRegexInput, 'g')
    const reBody = new RegExp(bodyRegexInput, 'g')
    const foundTixInTitle = reTitle.test(prTitle)
    const foundTixInBody = reBody.test(prBody)
    const foundNoTixInTitle = prTitle.startsWith(noTicketInput)
    const foundNoTixInBody = prBody.includes(noTicketInput)
    core.info(`found tix in title: ${foundTixInTitle}`)
    core.info(`found tix in body : ${foundTixInBody}`)
    core.info(`found no-tix in title: ${foundNoTixInTitle}`)
    core.info(`found no-tix in body : ${foundNoTixInBody}`)

    if (foundTixInTitle && foundTixInBody) {
      core.info('success, found tix in both title and body')
      return
    }
    if (foundNoTixInTitle && foundNoTixInBody) {
      core.info(`success, found ${noTicketInput} in both title and body`)
      return
    }
    if (isCI) {
      core.info(`ci mode detected. not returning failure`)
      return
    }

    core.setFailed('missing JIRA ticket in both PR title AND body')
    if (makePrComment === true) {
      createPRComment(client)
    }
  } catch (error) {
    core.setFailed(error.message)
  }
}

async function createPRComment(client: InstanceType<typeof GitHub>): Promise<void> {
  try {
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
