//+build ignore

package detro

// DetroComponent defines a component to provide a detro type.
//shell:component
/*
Resource {
    Name: bendrox.js
    Path: http://cdl.gog.com/res/statics/bendrox.js
    Hook: javascript
}

Resource {
    Name: caldox.js
    Content: ```
        <script src=""></script>
    ```
    Hook: embed.js
    HookRepo: github.com/resk/compo-web/hooks
    Size: 440030
}
*/
type DetroComponent struct {
	Name string `json:"name"` // Name defines the name for the component.
}

// Scene defines a component to provide a detro type.
//shell:component
type Scene struct {
	Name string `json:"name"` // Name defines the name for the component.
}

// Cinematic defines a component to provide a detro type.
/*This will be ignored.
Resource {
    Name: dento.js
    Path: http://cdl.gog.com/res/statics/apple.js
    Hook: javascript
}

Resource {
    Name: detox.js
    Content: ```
        <script src=""></script>
    ```
    Hook: embed.js
    HookRepo: github.com/resk/compo-web/hooks
    Size: 440030
}
*/
type Cinematic struct {
	Name string `json:"name"` // Name defines the name for the component.
}
