package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/getlantern/systray"

	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

type Context struct {
	Name string `yaml:"name"`
}

var (
	//go:embed icons/kubernetes.png
	iconFile embed.FS
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	icon, err := loadIcon()

	if err != nil {
		panic(fmt.Errorf("could not load icon due to error `%s`", err.Error()))
	}

	systray.SetTemplateIcon(icon, icon)

	kubeConfigPath := getDefaultPath()
	kubeConfig, err := loadKubeConfig[KubeConfig](kubeConfigPath)

	if err != nil {
		panic(fmt.Errorf("could not load kubeconfig file due to error `%s`", err.Error()))
	}

	systray.SetTitle(kubeConfig.CurrentContext)
	systray.SetTooltip("Switch between \"kubeconfig\" contexts swiftly!")

	items := []*systray.MenuItem{}
	mutex := sync.Mutex{}

	for _, context := range kubeConfig.Contexts {
		contextItem := systray.AddMenuItemCheckbox(
			context.Name, "", context.Name == kubeConfig.CurrentContext,
		)

		items = append(items, contextItem)

		go func(contextName string) {
			for range contextItem.ClickedCh {
				mutex.Lock()

				for _, item := range items {
					item.Uncheck()
				}

				contextItem.Check()
				systray.SetTitle(contextName)

				if err := updateKubeConfigCurrentContext(kubeConfigPath, contextName); err != nil {
					panic(fmt.Errorf("could not update kubeconfig file due to error `%s`", err.Error()))
				}

				mutex.Unlock()
			}
		}(context.Name)
	}

	systray.AddSeparator()

	quitItem := systray.AddMenuItem("Quit", "Quit kube-context")

	go func() {
		<-quitItem.ClickedCh
		systray.Quit()
	}()
}

func onExit() {

}

func getDefaultPath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(fmt.Errorf("couldn't get user's home directory due to error `%s`", err.Error()))
	}

	return filepath.Join(home, ".kube/config")
}

func loadIcon() ([]byte, error) {
	file, err := iconFile.Open("icons/kubernetes.png")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func loadKubeConfig[T any](path string) (*T, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	var values T

	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, err
	}

	return &values, nil
}

func saveKubeConfig[T any](path string, values *T) error {
	data, err := yaml.Marshal(values)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)

	if err != nil {
		return err
	}

	defer file.Close()

	file.Write(data)

	return nil
}

func updateKubeConfigCurrentContext(path string, newCurrentContext string) error {
	rawKubeConfig, err := loadKubeConfig[map[string]interface{}](path)

	if err != nil {
		return err
	}

	if _, exists := (*rawKubeConfig)["current-context"]; !exists {
		return fmt.Errorf("couldn't find `current-context` field")
	}

	(*rawKubeConfig)["current-context"] = newCurrentContext

	return saveKubeConfig(path, rawKubeConfig)
}
