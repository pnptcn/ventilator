package service

type PageInfo struct {
    URL         string            `json:"url"`
    Title       string            `json:"title"`
    Description string            `json:"description,omitempty"`
    Keywords    []string          `json:"keywords,omitempty"`
    Author      string            `json:"author,omitempty"`
    Trackers    map[string]interface{} `json:"trackers,omitempty"`
}

type PageContext struct {
    LocalStorage   map[string]interface{} `json:"localStorage"`
    SessionStorage map[string]interface{} `json:"sessionStorage"`
    Cookies        map[string]string      `json:"cookies"`
}

type PageContent struct {
    Raw    string   `json:"raw"`
    Chunks []string `json:"chunks"`
}

type Page struct {
    Info    PageInfo    `json:"info"`
    Context PageContext `json:"context"`
    Content PageContent `json:"content"`
}

type Profile struct {
    Page Page `json:"page"`
}
