/** @type {import("eslint").Linter.Config} */
module.exports = {
	extends: ['@repo/eslint-config'],
	root: true,
	env: { browser: true, es2020: true },
	plugins: ['react-refresh'],
	rules: {
		'react-refresh/only-export-components': ['warn', { allowConstantExport: true }]
	}
};
