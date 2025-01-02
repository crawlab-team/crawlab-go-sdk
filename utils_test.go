package crawlab

import (
	"bytes"
	"encoding/json"
	"github.com/crawlab-team/crawlab-go-sdk/constants"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveItem(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test single item
	item := map[string]any{
		"name": "test",
		"age":  25,
	}

	// Capture stdout
	output := captureOutput(func() {
		SaveItem(item)
	})

	// Verify output
	var msg entity.IPCMessage
	err := json.Unmarshal([]byte(output), &msg)
	require.NoError(err, "Should unmarshal JSON without error")

	assert.Equal(constants.IPCMessageTypeData, msg.Type, "Message type should be 'data'")
	assert.True(msg.IPC, "IPC flag should be true")

	items, ok := msg.Payload.([]any)
	assert.True(ok, "Payload should be an array")
	require.Len(items, 1, "Should have exactly 1 item")

	// Verify item contents
	firstItem, ok := items[0].(map[string]any)
	assert.True(ok, "Item should be a map")
	assert.Equal("test", firstItem["name"], "Item name should match")
	assert.Equal(float64(25), firstItem["age"], "Item age should match")
}

func TestSaveItems(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test multiple items
	items := []map[string]any{
		{
			"name": "test1",
			"age":  25,
		},
		{
			"name": "test2",
			"age":  30,
		},
	}

	// Capture stdout
	output := captureOutput(func() {
		SaveItems(items)
	})

	// Verify output
	var msg entity.IPCMessage
	err := json.Unmarshal([]byte(output), &msg)
	assert.NoError(err, "Should unmarshal JSON without error")

	assert.Equal(constants.IPCMessageTypeData, msg.Type, "Message type should be 'data'")
	assert.True(msg.IPC, "IPC flag should be true")

	payloadItems, ok := msg.Payload.([]any)
	assert.True(ok, "Payload should be an array")
	require.Len(payloadItems, 2, "Should have exactly 2 items")

	// Verify items contents
	firstItem, ok := payloadItems[0].(map[string]any)
	assert.True(ok, "First item should be a map")
	assert.Equal("test1", firstItem["name"], "First item name should match")
	assert.Equal(float64(25), firstItem["age"], "First item age should match")

	secondItem, ok := payloadItems[1].(map[string]any)
	assert.True(ok, "Second item should be a map")
	assert.Equal("test2", secondItem["name"], "Second item name should match")
	assert.Equal(float64(30), secondItem["age"], "Second item age should match")
}

// Helper function to capture stdout
func captureOutput(f func()) string {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	f()

	// Close the write end of the pipe to flush it
	w.Close()

	// Read the output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Restore original stdout
	os.Stdout = oldStdout

	return buf.String()
}
