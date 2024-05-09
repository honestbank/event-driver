/**
 * DO NOT CHANGE. This file is being managed from a central repository
 * To know more simply visit https://github.com/honestbank/.github/blob/main/docs/about.md
 */

module.exports = {
    branches: [{name: 'main'}],
    plugins: [
        ["@semantic-release/commit-analyzer", {
            "preset": "angular",
            "releaseRules": [
                {type: 'feat', release: 'minor'},
                {type: 'fix', release: 'patch'},
                {type: 'perf', release: 'patch'},
                {type: 'docs', release: 'patch'},
                {type: 'refactor', release: 'patch'},
                {type: 'style', release: 'patch'},
                {type: 'ci', release: 'patch'},
                {type: 'chore', release: 'patch'}
            ]
        }],
        "@semantic-release/release-notes-generator",
        "@semantic-release/github"
    ],
    verifyConditions: [
        "@semantic-release/github"
    ],
    prepare: [
    ],
    publish: [
        "@semantic-release/github"
    ],
};
