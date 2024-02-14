## Tasks

### generate

```sh
templ generate
```

### deploy

requires: generate
dir: cdk

```sh
cdk deploy
```

### deploy-hotswap

requires: generate
dir: cdk

```sh
cdk deploy --hotswap
```

### Project structure
- 'routes' are defined in the 'routes' directory
- 'routes' use services to get data that is passed down into templ 'components'. Each route exports a Handler. If there are sub-routes, they are imported and handled here.
- 'components' are defined in the 'components' directory
- Each route has a folder in the 'components' directory
- The 'components' directory also has a 'util' directory for shared components.
- Components do not have any application logic of their own. Components should not make API calls or manage application state.
- Components use View Models to reshape data for rendering.
- The 'eid' package is used only by components to generate unique HTML element IDs. The eid instance must be passed down through ALL components and never be reinstantiated outside the root main.go file.
