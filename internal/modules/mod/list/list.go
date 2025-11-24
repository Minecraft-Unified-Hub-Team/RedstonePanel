package list

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Minecraft-Unified-Hub-Team/RedstonePanel/internal/domain"
	"io"
)

type modDataRawTOML struct {
	DisplayName string `toml:"displayName"`
	ModID       string `toml:"modId"`
}

type modDataArrayTOML struct {
	Mods []modDataRawTOML `toml:"mods"`
}

func GetModDataFabric(file string) (domain.ModRef, error) {
	r, err := zip.OpenReader(file)
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

	data, err := io.ReadAll(rc)
	if err != nil {
		return domain.ModRef{}, fmt.Errorf("failed to read mod data at: %s", file)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return domain.ModRef{}, fmt.Errorf("JSON mod data is invalid at: %s", file)
	}

	name, ok := content["name"].(string)
	if !ok {
		return domain.ModRef{}, fmt.Errorf("failed to get mod name of: %s", file)
	}

	modid, ok := content["id"].(string)
	if !ok {
		return domain.ModRef{}, fmt.Errorf("failed to get mod ID of: %s", file)
	}

	return domain.ModRef{domain.ModID(modid), name, domain.MinecraftVersion(""), domain.URL("")}, nil
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

type ModList interface {
	List(ctx context.Context, srv domain.ServerID) ([]domain.ModRef, error)
}
