package htmlRender

import (
	"path"

	"github.com/gin-contrib/multitemplate"
)

const webDir = "web/templates"

func CreateMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", path.Join(webDir, "login.html"))
	r.AddFromFiles("config", path.Join(webDir, "layout.html"), path.Join(webDir, "config.html"))
	r.AddFromFiles("404", path.Join(webDir, "404.html"))
	r.AddFromFiles("500", path.Join(webDir, "500.html"))
	r.AddFromFiles("search", path.Join(webDir, "search.html"))
	r.AddFromFiles("main", path.Join(webDir, "layout.html"), path.Join(webDir, "main.html"))
	r.AddFromFiles("task", path.Join(webDir, "layout.html"), path.Join(webDir, "task.html"))
	r.AddFromFiles("plugin", path.Join(webDir, "layout.html"), path.Join(webDir, "plugin.html"))
	r.AddFromFiles("analysis", path.Join(webDir, "analysis.html"))
	return r
}
