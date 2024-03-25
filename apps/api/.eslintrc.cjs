/** @type {import("eslint").Linter.Config} */
module.exports = {
	extends: ['@repo/eslint-config'],
	root: true,
	env: { es2020: true }
};
