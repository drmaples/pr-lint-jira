import './sourcemap-register.cjs';/******/ /* webpack/runtime/compat */
/******/ 
/******/ if (typeof __nccwpck_require__ !== 'undefined') __nccwpck_require__.ab = new URL('.', import.meta.url).pathname.slice(import.meta.url.match(/^file:\/\/\/\w:/) ? 1 : 0, -1) + "/";
/******/ 
/************************************************************************/
var __webpack_exports__ = {};

/* eslint no-console: */
var __createBinding = (undefined && undefined.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (undefined && undefined.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (undefined && undefined.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const core = __importStar(require("@actions/core"));
const github = __importStar(require("@actions/github"));
const process = __importStar(require("process"));
async function run() {
    try {
        core.debug(JSON.stringify(github.context));
        const isCI = process.env.IS_CI === 'true';
        const token = core.getInput('token', { required: true });
        const makePrComment = core.getInput('make_pr_comment', { required: false }) === 'true';
        const titleRegexInput = core.getInput('title_regex', { required: false }) || '^\\[([A-Z]{2,}-\\d+)\\]';
        const bodyRegexInput = core.getInput('body_regex', { required: false }) || '\\[([A-Z]{2,}-\\d+)\\]';
        const noTicketInput = core.getInput('no_ticket', { required: false }) || '[no-ticket]';
        console.log(`input make_pr_comment: ${makePrComment}`);
        console.log(`input title_regex    : ${titleRegexInput}`);
        console.log(`input body_regex     : ${bodyRegexInput}`);
        console.log(`input no_ticket      : ${noTicketInput}`);
        const client = github.getOctokit(token);
        const prTitle = github.context?.payload?.pull_request?.title || '';
        const prBody = github.context?.payload?.pull_request?.body || '';
        const reTitle = new RegExp(titleRegexInput, 'g');
        const reBody = new RegExp(bodyRegexInput, 'g');
        const foundTixInTitle = reTitle.test(prTitle);
        const foundTixInBody = reBody.test(prBody);
        const foundNoTixInTitle = prTitle.startsWith(noTicketInput);
        const foundNoTixInBody = prBody.includes(noTicketInput);
        core.info(`found ticket in title   : ${foundTixInTitle}`);
        core.info(`found ticket in body    : ${foundTixInBody}`);
        core.info(`found no_ticket in title: ${foundNoTixInTitle}`);
        core.info(`found no_ticket in body : ${foundNoTixInBody}`);
        if (github.context?.eventName !== 'pull_request') {
            core.info('success, event is not a pull request');
            return;
        }
        if (foundTixInTitle && foundTixInBody) {
            core.info('success, found tix in both title and body');
            return;
        }
        if (foundNoTixInTitle && foundNoTixInBody) {
            core.info(`success, found ${noTicketInput} in both title and body`);
            return;
        }
        if (isCI) {
            core.info('ci mode detected. not returning failure');
            return;
        }
        core.setFailed('missing ticket in both PR title AND body');
        if (makePrComment === true) {
            createPRComment(client);
        }
    }
    catch (error) {
        core.setFailed(error.message);
    }
}
async function createPRComment(client) {
    try {
        // https://octokit.github.io/rest.js/v18#pulls-create-review-comment
        await client.rest.issues.createComment({
            owner: github.context.issue.owner,
            repo: github.context.issue.repo,
            issue_number: github.context.issue.number,
            body: 'PR title AND body does not contain a reference to a ticket.'
        });
    }
    catch (error) {
        core.error(`Failed to create PR comment: ${error}`);
    }
}
run();


//# sourceMappingURL=index.js.map