package crawlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/crawlab-team/crawlab-go-sdk/constants"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollyOnHTMLOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name          string
		html          string
		goqueryString string
		expectedItem  map[string]any
	}{
		{
			name:          "simple div extraction",
			html:          `<html><body><div class="item">Test Content</div></body></html>`,
			goqueryString: "div.item",
			expectedItem:  map[string]any{"content": "Test Content"},
		},
		{
			name:          "div with attributes",
			html:          `<html><body><div class="item" data-id="123">Test Content</div></body></html>`,
			goqueryString: "div.item",
			expectedItem:  map[string]any{"content": "Test Content", "id": "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(tt.html))
			}))
			defer ts.Close()

			// Create collector
			c := colly.NewCollector()

			// Set up the HTML callback
			CollyOnHTMLOne(c, tt.goqueryString, func(e *colly.HTMLElement) map[string]any {
				if tt.name == "div with attributes" {
					return map[string]any{
						"content": e.Text,
						"id":      e.Attr("data-id"),
					}
				}
				return map[string]any{
					"content": e.Text,
				}
			})

			// Capture stdout and visit the test server
			output := captureOutput(func() {
				err := c.Visit(ts.URL)
				require.NoError(err)
			})

			// Verify output format
			var msg entity.IPCMessage
			err := json.Unmarshal([]byte(output), &msg)
			require.NoError(err, "Should unmarshal JSON without error")

			assert.Equal(constants.IPCMessageTypeData, msg.Type, "Message type should be 'data'")
			assert.True(msg.IPC, "IPC flag should be true")

			// Verify items
			items, ok := msg.Payload.([]any)
			assert.True(ok, "Payload should be an array")
			require.Len(items, 1, "Should have exactly 1 item")

			// Convert and compare item
			actualItem, ok := items[0].(map[string]any)
			assert.True(ok, "Item should be a map")
			for k, v := range tt.expectedItem {
				assert.Equal(v, actualItem[k], "Item %s should match", k)
			}
		})
	}
}

func TestCollyOnHTMLMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name          string
		html          string
		goqueryString string
		expectedItems []map[string]any
	}{
		{
			name:          "single item extraction",
			html:          `<html><body><div class="item">Test Content</div></body></html>`,
			goqueryString: "div.item",
			expectedItems: []map[string]any{
				{"content": "Test Content"},
			},
		},
		{
			name: "multiple items extraction",
			html: `<html><body>
				<div class="item">First Item</div>
				<div class="item">Second Item</div>
			</body></html>`,
			goqueryString: "div.item",
			expectedItems: []map[string]any{
				{"content": "First Item"},
				{"content": "Second Item"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(tt.html))
			}))
			defer ts.Close()

			// Create collector
			c := colly.NewCollector()

			// Set up the HTML callback
			CollyOnHTMLMany(c, tt.goqueryString, func(e *colly.HTMLElement) []map[string]any {
				return []map[string]any{
					{"content": e.Text},
				}
			})

			// Capture stdout and visit the test server
			output := captureOutput(func() {
				err := c.Visit(ts.URL)
				require.NoError(err)
			})

			// Verify output format
			outputLines := strings.Split(strings.TrimSpace(output), "\n")
			assert.Len(outputLines, len(tt.expectedItems), "Should have expected number of output lines")
			for i, line := range outputLines {
				var msg entity.IPCMessage
				err := json.Unmarshal([]byte(line), &msg)
				require.NoError(err, "Should unmarshal JSON without error")

				assert.Equal(constants.IPCMessageTypeData, msg.Type, "Message type should be 'data'")
				assert.True(msg.IPC, "IPC flag should be true")

				// Verify item
				expectedItem := tt.expectedItems[i]
				payloadItems, ok := msg.Payload.([]interface{})
				require.True(ok)
				actualItem, ok := payloadItems[0].(map[string]any)
				require.True(ok)
				assert.True(ok, "Payload should be a map")
				for k, v := range expectedItem {
					assert.Equal(v, actualItem[k], "Item %s should match", k)
				}
			}
		})
	}
}
