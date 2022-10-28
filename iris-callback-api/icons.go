package iriscallbackapi

type iconsList struct {
	Warn       string
	Success    string
	SuccessOff string
	Notice     string
	Info       string
	Danger     string
	Comment    string
	Config     string
	Catalog    string
	Stats      string
}

var (
	Icons = iconsList{
		Warn:       "⚠️ ",
		Success:    "✅ ",
		SuccessOff: "❎ ",
		Notice:     "📝 ",
		Info:       "🗓 ",
		Danger:     "📛 ",
		Comment:    "💬 ",
		Config:     "⚙️ ",
		Catalog:    "🗂 ",
		Stats:      "📊 ",
	}
)
