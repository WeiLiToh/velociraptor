package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Velocidex/ordereddict"
	vql_subsystem "www.velocidex.com/golang/velociraptor/vql"
	"www.velocidex.com/golang/vfilter"
	"www.velocidex.com/golang/vfilter/types"
)

type OllamaPlugin struct{}

func (self *OllamaPlugin) Info(scope types.Scope, type_map *vfilter.TypeMap) *vfilter.PluginInfo {
	return &vfilter.PluginInfo{
		Name: "ollama",
		Doc:  "Send VQL query results to an Ollama LLM (e.g., Qwen2.5) for automated triage/analysis.",
	}
}

func (self *OllamaPlugin) Call(ctx context.Context, scope types.Scope, args *ordereddict.Dict) <-chan types.Row {
	output := make(chan types.Row)

	go func() {
		defer close(output)

		inputVal, exists := args.Get("input")
		if !exists {
			row := ordereddict.NewDict()
			row.Set("error", "Missing required parameter: input")
			output <- row
			return
		}

		inputRows, ok := inputVal.([]interface{})
		if !ok {
			row := ordereddict.NewDict()
			row.Set("error", "Invalid input format (expected a list)")
			output <- row
			return
		}

		var rows []map[string]interface{}
		for _, r := range inputRows {
			if rowMap, ok := r.(map[string]interface{}); ok {
				rows = append(rows, rowMap)
			}
		}

		model := "qwen:2.5"
		if m, ok := args.Get("model"); ok {
			model = fmt.Sprintf("%v", m)
		}

		prompt := fmt.Sprintf(`You are a digital forensic analyst. Analyze the following VQL output for suspicious behavior:

%v`, toPrettyJSON(rows))

		req := map[string]interface{}{
			"model":  model,
			"prompt": prompt,
			"stream": false,
		}
		body, _ := json.Marshal(req)

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
		if err != nil {
			row := ordereddict.NewDict()
			row.Set("error", fmt.Sprintf("Ollama HTTP error: %v", err))
			output <- row
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		var llmResp map[string]interface{}
		if err := json.Unmarshal(respBody, &llmResp); err != nil {
			row := ordereddict.NewDict()
			row.Set("error", "Failed to parse LLM response")
			output <- row
			return
		}

		row := ordereddict.NewDict()
		row.Set("llm_response", llmResp["response"])
		row.Set("rows_input", len(rows))
		row.Set("model_used", model)
		output <- row
	}()

	return output
}

func toPrettyJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}
	return string(b)
}

// Register plug-in
func init() {
	vql_subsystem.RegisterPlugin(&OllamaPlugin{})
}
