```
brew install argocd
argocd app create root --upsert \
    --dest-namespace argocd \
    --dest-server https://kubernetes.default.svc \
    --repo https://github.com/brunoluiz/go-lab.git \
    --path infrastructure/argocd/root
argocd app sync root --prune
```
