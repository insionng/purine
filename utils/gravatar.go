package utils

func GravatarLink(email string) string {
	return "https://gravatar.com/avatar/" + Md5String(email)
}
