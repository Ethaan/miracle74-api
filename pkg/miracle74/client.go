package miracle74

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"golang.org/x/net/html"
)

const (
	baseURL        = "https://miracle74.com"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}
}

func (c *Client) ScrapeCharacter(name string) (*types.Character, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("subtopic", "characters")
	q.Set("name", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Miracle74-API/0.1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	character, err := c.parseCharacterHTML(body, name)
	if err != nil {
		return nil, fmt.Errorf("failed to parse character data: %w", err)
	}

	return character, nil
}

func (c *Client) parseCharacterHTML(htmlContent []byte, name string) (*types.Character, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	character, err := parseCharacterData(doc)
	if err != nil {
		fmt.Printf("DEBUG: Parser error: %v\n", err)
		return nil, fmt.Errorf("failed to extract character data: %w", err)
	}

	fmt.Printf("DEBUG: Parsed character: %+v\n", character)
	return character, nil
}

func (c *Client) ScrapePowerGamers(includeAll bool) ([]types.PowerGamer, error) {
	var allPowerGamers []types.PowerGamer

	maxPages := 1
	if includeAll {
		maxPages = 10
	}

	for page := 1; page <= maxPages; page++ {
		fmt.Printf("Scraping power gamers page %d...\n", page)

		u, err := url.Parse(c.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}

		q := u.Query()
		q.Set("subtopic", "powergamers")
		q.Set("list", "today")
		q.Set("page", fmt.Sprintf("%d", page))
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for page %d: %w", page, err)
		}

		req.Header.Set("User-Agent", "Miracle74-API/0.1.0")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code for page %d: %d", page, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for page %d: %w", page, err)
		}

		powerGamers, err := c.parsePowerGamersHTML(body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse power gamers data from page %d: %w", page, err)
		}

		allPowerGamers = append(allPowerGamers, powerGamers...)

		if page < maxPages {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("Successfully scraped %d power gamers from %d page(s)\n", len(allPowerGamers), maxPages)
	return allPowerGamers, nil
}

func (c *Client) parsePowerGamersHTML(htmlContent []byte) ([]types.PowerGamer, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	powerGamers, err := parsePowerGamersData(doc)
	if err != nil {
		fmt.Printf("DEBUG: Parser error: %v\n", err)
		return nil, fmt.Errorf("failed to extract power gamers data: %w", err)
	}

	fmt.Printf("DEBUG: Parsed %d power gamers from page\n", len(powerGamers))
	return powerGamers, nil
}

func (c *Client) ScrapeInsomniacs(includeAll bool) ([]types.Insomniac, error) {
	var allInsomniacs []types.Insomniac

	maxPages := 1
	if includeAll {
		maxPages = 10
	}

	for page := 1; page <= maxPages; page++ {
		fmt.Printf("Scraping insomniacs page %d...\n", page)

		u, err := url.Parse(c.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}

		q := u.Query()
		q.Set("subtopic", "insomniacs")
		q.Set("page", fmt.Sprintf("%d", page))
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for page %d: %w", page, err)
		}

		req.Header.Set("User-Agent", "Miracle74-API/0.1.0")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code for page %d: %d", page, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for page %d: %w", page, err)
		}

		insomniacs, err := c.parseInsomniacsHTML(body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse insomniacs data from page %d: %w", page, err)
		}

		allInsomniacs = append(allInsomniacs, insomniacs...)

		if page < maxPages {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("Successfully scraped %d insomniacs from %d page(s)\n", len(allInsomniacs), maxPages)
	return allInsomniacs, nil
}

func (c *Client) parseInsomniacsHTML(htmlContent []byte) ([]types.Insomniac, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	insomniacs, err := parseInsomniacsData(doc)
	if err != nil {
		fmt.Printf("DEBUG: Parser error: %v\n", err)
		return nil, fmt.Errorf("failed to extract insomniacs data: %w", err)
	}

	fmt.Printf("DEBUG: Parsed %d insomniacs\n", len(insomniacs))
	return insomniacs, nil
}

func (c *Client) ScrapeGuild(guildID int) (*types.Guild, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("subtopic", "guilds")
	q.Set("action", "show")
	q.Set("guild", fmt.Sprintf("%d", guildID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Miracle74-API/0.1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	guild, err := c.parseGuildHTML(body, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guild data: %w", err)
	}

	return guild, nil
}

func (c *Client) parseGuildHTML(htmlContent []byte, guildID int) (*types.Guild, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	guild, err := parseGuildData(doc, guildID)
	if err != nil {
		fmt.Printf("DEBUG: Parser error: %v\n", err)
		return nil, fmt.Errorf("failed to extract guild data: %w", err)
	}

	fmt.Printf("DEBUG: Parsed guild with %d members\n", len(guild.Members))
	return guild, nil
}

// func (c *Client) saveHTMLForDebug(htmlContent []byte, name string) error {
// 	publicDir := "public"

// 	if err := os.MkdirAll(publicDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create public directory: %w", err)
// 	}

// 	filename := filepath.Join(publicDir, fmt.Sprintf("%s.html", name))
// 	if err := os.WriteFile(filename, htmlContent, 0644); err != nil {
// 		return fmt.Errorf("failed to write HTML file: %w", err)
// 	}

// 	fmt.Printf("Saved HTML to %s\n", filename)
// 	return nil
// }
