package cli

import (
    "strings"
)

func TextSplitCSV(text string) []string {
    if text == "" {
        return nil
    }
    
    result := strings.Split(text,",")
    for i, _ := range result {
        result[i]=strings.TrimSpace(result[i])
    }
    return result
}
