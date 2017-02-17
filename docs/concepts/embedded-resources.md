Embedded Resources
==================

In using GU, there exists at times the need to use external resources. Understanding this need, Gu provides the ability to list resources for components which will be loaded on startup and based on when this components will be called. It allows the inclusion of different resources (e.g CSS, Javascript, Images), which are then  installed by custom hooks into the virtual DOM.

Internally GU uses a two stage process. The First stage involves parsing the intended package for Resource meta-data declarations and produces a `manifest.json` file. This file will be automatically loaded on startup and hence requires the developer to expose the respective path of the generated manifest file to be accessible on the server. The Second stage is when the manifest file is parsed to retrieve all resources and loads those which are required by the components of a view or based if it's declared as a global resource.

Additionally, the GU parser will search through all import paths to find additional resource declarations to be included for the calling package. Usually you only ever need to generate the `manifest.json` file for the package which will be executing. All resources declared by the application and it's imports will be included within that `manifest.json` file and will be loaded accordingly. This allows alot of simplication and provides a single source of truth for embeddable resources.

*GU provides a CLI tool that is instaleld when you `go get`  the GU package. It helps in generating the manifest.json file and also optionally creates a virtual file system which can be loaded as a package for single binary deployments.*

Declaring Resources
-------------------

Declaring resources to be embedded along with a component or package is easy. Gu uses the meta-data declarations (in the manifest of the golang code ?)  which will be scanned and pulled accordingly. Gu provides a set of fields usable when declaring a resource, as shown below.

```go
Resource {
  Name: string              // Custom name of resource which it will be accessed under. (REQUIRED)

  Path: string              // Custom path of resource if it's from a local or remote file on a CDN/Other endpoint. (OPTIONAL IF CONTENT IS PROVIDED)

  Localize: bool            // Declares whether the resource should be pulled from endpoint path and copied into manifests file. (REQUIRED)

  Relations: string         // A coma separated listing of components which use this resource. Helps to avoid duplications. (OPTIONAL)

  Hook: string              // Name of the hook which will install resource. (REQUIRED)

  Content: string           // Data of the resource if no path is provided. (REQUIRED If No Path exists)

  Size: int                 // Size to use for resource when serving. (OPTIONAL, DEFAULT: Automatically set when resource is localized)

  Init: bool                // Declares whether resource should be installed when component is initialized. (OPTIONAL, DEFAULT: true)

  Global: bool              // Declares whether resource is global and should be registered and accessible through the Gu resource registry. (OPTIONAL, DEFAULT: false)

  Encode: bool              // Declares whether resource should be base64 encoded when retrieved from endpoint path. (OPTIONAL, DEFAULT: true)

  Base64Padding: bool       // Declares if resource should be encoded with base64 padding or non padding. (REQUIRED)

  Encoded64: bool           // Declares if embedded content in Content field is already base64 encoded. (OPTIONAL, DEFAULT: false)

  Remote: bool              // Declares that resource is remote even if path provided is local and will be provided by serve. (OPTIONAL, DEFAULT: false)
}
```

The above fields define the behaviour and means by which a embedded resource should be processed and accessed by the Gu parser.

When the GU parser finds field names which do not match this giving set (not sure what "giving set" means ??) , then these are extracted into a map as meta-details, which can then be used by the hooks as implemented to achieve desired resources or as decorations for a more detailed resource.

-	Declaring Global Resource Generally when there exists resources which should be included on all views regardless of content, then the global resource declaration strategy should be used has it provides the capability to achieve this.

By simply declaring the resource meta-data at the package declaration, the parser will mark these resources as global and will be used on every rendering of the application.

```go
// Package component contains useful components built with Gu.
//
//shell:component:global
//
// Resource {
//     Name: detos.font.js
//     Path: https://fonts.googleapis.com/css?family=Lato|Open+Sans|Roboto
//     Hook: embedded-css
//		 Localize: true
// }
//
package components

// Menu component.
type Menu struct{}
```

Global resources are marked by the `shell:component:global` marker and this is always required and must be declared first and separate from any package level comments. This marker will be used by the GU parser to identify resource declarations which are required globally.

**This declaration can not be used any where else and must be declared immediately after a package comments not after.**

-	Declaring Component based Resources When declaring resources specific to the existence and initialization of specific components. The resource need to be declared after the said component's commentary.

Gu uses a `shell:component` marker to site if a giving type declares resources which are required to be embedded. Also the Gu parser is wise enough to included meta details corresponding to a components internal component and declared types.

These information then will be used to loaded addition resources which relate to those types when the giving component is being initialized.

```go
// Package component contains useful components built with Gu.
package components

import "github.com/gupa/components"

// Menu component.
//
//shell:component
//
// Resource {
//     Name: detox.js
//     Path: https://fonts.googleapis.com/detox.js
//		 Localize: false
//		 Relations: BobComponent, HUDComponent
//     Hook: js
// }
//
// Resource {
// 		Version: 1.4.1
// 		Pkg: DentusVuz
//     Name: hul_hub.js
//     Content: ```
//         <script src="">
// 					var mo = func() Resource {
// 						return ```block Resource {
// 							name: 'bock',
// 						}```
// 					}
// 				</script>
//     ```
//     Hook: embed.js
//     HookRepo: github.com/resk/compo-web/hooks
//     Size: 440030
// }
//
type Menu struct {
	List components.List
}
```

With the above approach. It becomes generally easy to quickly declare and have resources quickly embedded into packages for components.

Generate `manifest.json` file
-----------------------------

*As adviced, the `manifest.json` file should be only generated for the package with the `main` function which will be the means by which the application is lunched.*

Gu provides a CLI tool with the package once installed with `go get`:

```bash
go get -u github.com/gu-io/gu/...
```

The Gu CLI then is installed and then can be accessible through the terminal to generate the manifest file or a package which can be included on the server to provide a virtual file system for the generated resources.

-	Generating the `manifest.json` file.

```bash
gu generate --indir ./apps --outdir ./
```

Where the fields: --indir: Defines the directories to search in for resource embedding meta declarations --outdir: Defines the directory where the `manifest.json` file should be stored.

-	Generating the `manifest.json` file.

```bash
gu generate-vfs --indir ./apps --outdir ./ -pkg manifest-vfs
```

Where the fields: --indir: Defines the directories to search in for resource embedding meta declarations --outdir: Defines the directory where the new virtual file system package should be stored. --pkg: Defines the name of the package which is reflective on the name of the package directory (Default: "manifests").
