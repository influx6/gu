package main

import (
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/examples/greeter/components"
	"github.com/gu-io/gu/trees/elems"
)

func main() {
	app := gu.App(AppSettings)

	index := app.View(gu.ViewAttr{
		Name:  "View.Greeter",
		Route: "/*",
		Base: elems.Parse(`
			<div class="greeter-view view wrapper">
				<h1 class="view-header">Greeter App</h1>

				<div class="greeter-app" id="greeter-app-component">
				</div>
			</div>
		`, elems.CSS(`
				html *{
					padding: 0;
					margin: 0;
					font-size:16px;
					font-size: 100%;
					box-sizing: border-box;
				}

				html {
					width: 100%;
					height: 100%;
				}

				body{
					width: 100%;
					height: 100%;
					font-family: "Lato", helvetica, sans-serif;
					background: url("assets/galaxy3.jpg") no-repeat;
					background-size: cover;
				}

				&{
					color: #fff;
					width: 100%;
					padding: 10px;
					min-height: 100%;
					margin: 0px auto;
					background: rgba(0,0,0,0.4);
				}

				& h1{
					text-align: center;
					font-size: 2.5em;
					margin: 40px auto;
				}

				& .greeter-app {
					width: 90%;
					height: auto;
					margin: 30px auto;
					padding-top: 100px;
					text-align: center;
				}


				& .greeter-app .receiver{
					font-size: 1.7em;
				}

				& .greeter-app .receiver input{
					color: #fff !important;
				}
		`, nil)),
	})

	index.Component(gu.ComponentAttr{
		Route:  "/*",
		Target: "#greeter-app-component",
		Base:   components.NewGreeter(),
	})
}
