/** @type {import("eslint").Linter.Config} */
module.exports = {
	extends: ['@repo/eslint-config'],
	ignorePatterns: ['apps/**', 'packages/**'],
	root: true,
	parser: '@typescript-eslint/parser',
	parserOptions: {
		project: true
	}
};
