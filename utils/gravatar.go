package utils

// create gravatar link
func GravatarLink(email string) string {
	return "https://gravatar.com/avatar/" + Md5String(email)
}
