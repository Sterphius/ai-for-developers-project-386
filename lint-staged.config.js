export default {
  '*.go': () => 'go vet ./server/...',
  '*.{ts,tsx}': () => 'npm run lint --prefix web',
};
