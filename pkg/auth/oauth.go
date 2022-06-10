package auth

type OAuth struct {
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

func (oauth *OAuth) New() {

}