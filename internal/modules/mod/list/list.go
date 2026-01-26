package list

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Minecraft-Unified-Hub-Team/RedstonePanel/internal/domain"
	"io"
)

type modDataRaw struct {
	DisplayName string `toml:"displayName" json:"name"`
	ModID       string `toml:"modId" json:"id"`
}

type modDataArrayTOML struct {
	Mods []modDataRaw `toml:"mods"`
}

type zipManager interface {
	OpenReader(file string) (*zip.ReadCloser, error)
}

type Lister struct {
	zipManager zipManager
}

func NewLister(zipManager zipManager) *Lister {
	return &Lister{
		zipManager: zipManager,
	}
}

func (l *Lister) GetModDataFabric(file string) (domain.ModRef, error) {
	r, err := l.zipManager.OpenReader(file)
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to access mod file: %s", file)
	}
	defer r.Close()

	var targetJson *zip.File
	for _, f := range r.File {
		if f.Name == "fabric.mod.json" {
			targetJson = f
			break
		}
	}
	if targetJson == nil {
		return domain.ModRef{}, fmt.Errorf("failed to locate mod data at: %s", file)
	}

	rc, err := targetJson.Open()
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to access mod data at: %s", file)
	}

	data0, err := io.ReadAll(rc)
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to read mod data at: %s", file)
	}

	var data modDataRaw
	if err := json.Unmarshal(data0, &data); err != nil {
		return domain.ModRef{}, fmt.Errorf("JSON mod data is invalid at: %s", file)
	}

	return domain.ModRef{domain.ModID(data.ModID), data.DisplayName, domain.MinecraftVersion(""), domain.URL("")}, nil
}

func GetModDataForge(file string) (domain.ModRef, error) {
	r, err := zip.OpenReader(file)
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to access mod file: %s", file)
	}
	defer r.Close()

	var targetToml *zip.File
	for _, f := range r.File {
		if f.Name == "META-INF/mods.toml" {
			targetToml = f
			break
		}
	}
	if targetToml == nil {
		return domain.ModRef{}, fmt.Errorf("failed to locate mod data at: %s", file)
	}

	rc, err := targetToml.Open()
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to access TOML mod data at: %s", file)
	}

	data0, err := io.ReadAll(rc)
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to read TOML mod data at: %s", file)
	}

	var dataArray modDataArrayTOML
	if err := toml.Unmarshal(data0, &dataArray); err != nil {
		return domain.ModRef{}, fmt.Errorf("TOML mod data is invalid at: %s", file)
	}

	return domain.ModRef{domain.ModID(dataArray.Mods[0].ModID), dataArray.Mods[0].DisplayName,
		domain.MinecraftVersion(""), domain.URL("")}, nil
}
