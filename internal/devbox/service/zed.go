package service

func (o *Operator) ZedURL(projectName string) (string, error) {
	project, err := o.Store.GetProject(projectName)
	if err != nil {
		return "", err
	}
	host := envOrDefault("DEVBOX_HOST_SSH_NAME", "thebox")
	return "zed ssh://" + host + project.WorkspacePath, nil
}
