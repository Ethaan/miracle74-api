package miracle74

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethaan/miracle74-api/internal/types"
	"golang.org/x/net/html"
)

func parseCharacterData(doc *html.Node) (*types.Character, error) {
	character := &types.Character{}

	table := findCharacterTable(doc)
	if table == nil {
		return nil, fmt.Errorf("character information table not found")
	}

	parseCharacterInfo(table, character)

	deathsTable := findDeathsTable(doc)
	if deathsTable != nil {
		character.Deaths = parseDeaths(deathsTable)
	}

	return character, nil
}

func findCharacterTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		if hasClass(n, "TableContent") && hasClass(n, "InnerBorder") {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findCharacterTable(c); result != nil {
			return result
		}
	}

	return nil
}

func findDeathsTable(n *html.Node) *html.Node {
	if isDeathsSection(n) {
		return findNextTable(n.Parent)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findDeathsTable(c); result != nil {
			return result
		}
	}

	return nil
}

func isDeathsSection(n *html.Node) bool {
	text := getTextContent(n)
	return strings.Contains(text, "Character Deaths")
}

func findNextTable(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if result := findFirstTable(s); result != nil {
			return result
		}
	}

	return findNextTable(n.Parent)
}

func findFirstTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findFirstTable(c); result != nil {
			return result
		}
	}

	return nil
}

func parseCharacterInfo(table *html.Node, character *types.Character) {
	rows := findAllTRs(table)

	for _, row := range rows {
		cells := findAllTDs(row)
		if len(cells) < 2 {
			continue
		}

		label := strings.TrimSpace(getTextContent(cells[0]))
		value := strings.TrimSpace(getTextContent(cells[1]))

		switch {
		case strings.Contains(label, "Name:"):
			character.Name = extractName(value)
			character.Country = extractCountry(cells[1])

		case strings.Contains(label, "Sex:"):
			character.Sex = value

		case strings.Contains(label, "Vocation:"):
			character.Vocation = value

		case strings.Contains(label, "Level:"):
			if level, err := strconv.Atoi(value); err == nil {
				character.Level = level
			}

		case strings.Contains(label, "Residence:"):
			character.Residence = value

		case strings.Contains(label, "Guild Membership:"):
			character.Guild, character.GuildRank, character.GuildURL = extractGuildInfo(cells[1])

		case strings.Contains(label, "Last login:"):
			if t, err := parseLastLogin(value); err == nil {
				character.LastLogin = &t
			}

		case strings.Contains(label, "Account") && strings.Contains(label, "Status:"):
			character.IsPremium = strings.Contains(value, "Premium")
		}
	}
}

func parseDeaths(table *html.Node) []types.Death {
	var deaths []types.Death
	rows := findAllTRs(table)

	for _, row := range rows {
		cells := findAllTDs(row)
		if len(cells) < 2 {
			continue
		}

		dateStr := strings.TrimSpace(getTextContent(cells[0]))
		deathInfo := strings.TrimSpace(getTextContent(cells[1]))

		level := extractDeathLevel(deathInfo)
		killedBy := extractKilledBy(deathInfo)

		if dateStr != "" && killedBy != "" {
			deaths = append(deaths, types.Death{
				Date:     dateStr,
				Level:    level,
				KilledBy: killedBy,
			})
		}
	}

	return deaths
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, class) {
			return true
		}
	}
	return false
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getTextContent(c)
	}
	return text
}

func findAllTRs(n *html.Node) []*html.Node {
	var rows []*html.Node

	if n.Type == html.ElementNode && n.Data == "tr" {
		rows = append(rows, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rows = append(rows, findAllTRs(c)...)
	}

	return rows
}

func findAllTDs(n *html.Node) []*html.Node {
	var cells []*html.Node

	if n.Type == html.ElementNode && n.Data == "td" {
		cells = append(cells, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cells = append(cells, findAllTDs(c)...)
	}

	return cells
}

func extractName(value string) string {
	lines := strings.Split(value, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return strings.TrimSpace(value)
}

func extractCountry(cell *html.Node) string {
	img := findFirstImg(cell)
	if img != nil {
		for _, attr := range img.Attr {
			if attr.Key == "src" && strings.Contains(attr.Val, "/images/flags/") {
				parts := strings.Split(attr.Val, "/")
				if len(parts) > 0 {
					filename := parts[len(parts)-1]
					return strings.TrimSuffix(filename, ".gif")
				}
			}
		}
	}
	return ""
}

func findFirstImg(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "img" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findFirstImg(c); result != nil {
			return result
		}
	}

	return nil
}

func extractGuildInfo(cell *html.Node) (guildName, guildRank, guildURL string) {
	var textBeforeLink string
	var linkNode *html.Node

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			linkNode = n
			return
		}
		if n.Type == html.TextNode && linkNode == nil {
			textBeforeLink += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(cell)

	if linkNode != nil {
		guildName = strings.TrimSpace(getTextContent(linkNode))

		for _, attr := range linkNode.Attr {
			if attr.Key == "href" {
				if strings.HasPrefix(attr.Val, "?") {
					guildURL = "https://miracle74.com/" + attr.Val
				} else {
					guildURL = attr.Val
				}
				break
			}
		}
	}

	textBeforeLink = strings.TrimSpace(textBeforeLink)
	if strings.Contains(textBeforeLink, " of the ") {
		parts := strings.Split(textBeforeLink, " of the ")
		if len(parts) > 0 {
			guildRank = strings.TrimSpace(parts[0])
		}
	} else {
		guildRank = textBeforeLink
	}

	return guildName, guildRank, guildURL
}

func parseLastLogin(value string) (time.Time, error) {
	layouts := []string{
		"2 January 2006, 3:04 pm",
		"2 January 2006, 3:04 am",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse last login: %s", value)
}

func extractDeathLevel(deathInfo string) int {
	re := regexp.MustCompile(`level (\d+)`)
	matches := re.FindStringSubmatch(deathInfo)
	if len(matches) > 1 {
		if level, err := strconv.Atoi(matches[1]); err == nil {
			return level
		}
	}
	return 0
}

func extractKilledBy(deathInfo string) string {
	parts := strings.Split(deathInfo, " by ")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func parsePowerGamersData(doc *html.Node) ([]types.PowerGamer, error) {
	table := findPowerGamersTable(doc)
	if table == nil {
		return nil, fmt.Errorf("power gamers table not found")
	}

	rows := findAllTRs(table)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found in power gamers table")
	}

	var powerGamers []types.PowerGamer

	for i, row := range rows {
		if i == 0 {
			text := getTextContent(row)
			if strings.Contains(text, "Rank") && strings.Contains(text, "Name") {
				continue
			}
		}

		cells := findAllTDs(row)
		if len(cells) < 5 {
			continue
		}

		rankStr := strings.TrimSpace(getTextContent(cells[0]))
		name := strings.TrimSpace(getTextContent(cells[1]))
		vocation := strings.TrimSpace(getTextContent(cells[2]))
		levelStr := strings.TrimSpace(getTextContent(cells[3]))
		todayStr := strings.TrimSpace(getTextContent(cells[4]))

		if rankStr == "" || name == "" {
			continue
		}

		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse rank '%s': %v\n", rankStr, err)
			continue
		}

		level, err := strconv.Atoi(levelStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse level '%s': %v\n", levelStr, err)
			continue
		}

		today, err := strconv.Atoi(todayStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse today '%s': %v\n", todayStr, err)
			continue
		}

		powerGamers = append(powerGamers, types.PowerGamer{
			Rank:     rank,
			Name:     name,
			Vocation: vocation,
			Level:    level,
			Today:    today,
		})
	}

	return powerGamers, nil
}

func findPowerGamersTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		if hasClass(n, "TableContent") && hasClass(n, "InnerBorder") {
			tbody := findTBody(n)
			if tbody != nil {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findPowerGamersTable(c); result != nil {
			return result
		}
	}

	return nil
}

func findTBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findTBody(c); result != nil {
			return result
		}
	}

	return nil
}

func parseInsomniacsData(doc *html.Node) ([]types.Insomniac, error) {
	table := findInsomniacsTable(doc)
	if table == nil {
		return nil, fmt.Errorf("insomniacs table not found")
	}

	rows := findAllTRs(table)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found in insomniacs table")
	}

	var insomniacs []types.Insomniac

	for i, row := range rows {
		if i == 0 {
			text := getTextContent(row)
			if strings.Contains(text, "Rank") || strings.Contains(text, "Name") {
				continue
			}
		}

		cells := findAllTDs(row)
		if len(cells) < 5 {
			continue
		}

		rankStr := strings.TrimSpace(getTextContent(cells[0]))
		nameCell := cells[1]
		vocation := strings.TrimSpace(getTextContent(cells[2]))
		levelStr := strings.TrimSpace(getTextContent(cells[3]))
		timeOnline := strings.TrimSpace(getTextContent(cells[4]))

		if rankStr == "" {
			continue
		}

		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse rank '%s': %v\n", rankStr, err)
			continue
		}

		name := extractNameFromLink(nameCell)
		if name == "" {
			fmt.Printf("Warning: failed to extract name from cell\n")
			continue
		}

		country := extractCountry(nameCell)

		level, err := strconv.Atoi(levelStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse level '%s': %v\n", levelStr, err)
			continue
		}

		insomniacs = append(insomniacs, types.Insomniac{
			Rank:       rank,
			Name:       name,
			Country:    country,
			Vocation:   vocation,
			Level:      level,
			TimeOnline: timeOnline,
		})
	}

	return insomniacs, nil
}

func findInsomniacsTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		if hasClass(n, "TableContent") && hasClass(n, "InnerBorder") {
			tbody := findTBody(n)
			if tbody != nil {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findInsomniacsTable(c); result != nil {
			return result
		}
	}

	return nil
}

func extractNameFromLink(cell *html.Node) string {
	link := findFirstLink(cell)
	if link != nil {
		return strings.TrimSpace(getTextContent(link))
	}
	return ""
}

func findFirstLink(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findFirstLink(c); result != nil {
			return result
		}
	}

	return nil
}

func parseGuildData(doc *html.Node, guildID int) (*types.Guild, error) {
	table := findGuildMembersTable(doc)
	if table == nil {
		return nil, fmt.Errorf("guild members table not found")
	}

	rows := findAllTRs(table)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found in guild members table")
	}

	var members []types.GuildMember

	for i, row := range rows {
		if i == 0 {
			text := getTextContent(row)
			if strings.Contains(text, "Rank") && strings.Contains(text, "Name") {
				continue
			}
		}

		cells := findAllTDs(row)
		if len(cells) < 5 {
			continue
		}

		rank := strings.TrimSpace(getTextContent(cells[0]))
		nameCell := cells[1]
		vocation := strings.TrimSpace(getTextContent(cells[2]))
		levelStr := strings.TrimSpace(getTextContent(cells[3]))
		statusCell := cells[4]

		if rank == "" {
			continue
		}

		name := extractNameFromLink(nameCell)
		if name == "" {
			fmt.Printf("Warning: failed to extract name from cell\n")
			continue
		}

		level, err := strconv.Atoi(levelStr)
		if err != nil {
			fmt.Printf("Warning: failed to parse level '%s': %v\n", levelStr, err)
			continue
		}

		status := extractGuildMemberStatus(statusCell)

		members = append(members, types.GuildMember{
			Rank:     rank,
			Name:     name,
			Vocation: vocation,
			Level:    level,
			Status:   status,
		})
	}

	return &types.Guild{
		GuildID: guildID,
		Members: members,
	}, nil
}

func findGuildMembersTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		if hasClass(n, "TableContent") && hasClass(n, "InnerBorder") {
			tbody := findTBody(n)
			if tbody != nil {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findGuildMembersTable(c); result != nil {
			return result
		}
	}

	return nil
}

func extractGuildMemberStatus(cell *html.Node) string {
	text := strings.TrimSpace(getTextContent(cell))

	if strings.Contains(text, "Online") || strings.Contains(text, "online") {
		return "Online"
	}

	return "Offline"
}
