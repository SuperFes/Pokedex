package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	expires time.Time
	data    []byte
	mu      sync.RWMutex
}

type CacheStore interface {
	Expired() bool
}

type Caches map[string]*Cache

func (c *Cache) Expired() bool {
	return time.Now().After(c.expires)
}

func (c *Caches) Set(name string, data []byte, duration time.Duration) {
	if *c == nil {
		*c = make(Caches)
	}

	if _, ok := (*c)[name]; ok {
		(*c)[name].mu.Lock()
		(*c)[name].expires = time.Now().Add(duration)
		(*c)[name].data = data
		(*c)[name].mu.Unlock()
	} else {
		(*c)[name] = &Cache{
			expires: time.Now().Add(duration),
			data:    data,
		}
	}

	for name, cache := range *c {
		if cache.Expired() {
			delete(*c, name)
		}
	}
}

func (c *Caches) Get(name string) ([]byte, bool) {
	if *c == nil {
		*c = make(Caches)
	}

	if cache, ok := (*c)[name]; ok {
		cache.mu.RLock()
		defer cache.mu.RUnlock()

		return cache.data, true
	}

	return nil, false
}

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type config struct {
	Commands map[string]cliCommand
	Map      struct {
		Next     string
		Previous string
	}
	Redis   *Caches
	Args    []string
	Pokemon map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	clean := make([]string, 0)

	parts := strings.Fields(strings.ToLower(strings.Trim(text, " ")))

	for _, part := range parts {
		if part != "" {
			clean = append(clean, part)
		}
	}

	return clean
}

func ApiGetLocations(url string, conf *config) (Location, error) {
	// GET https://pokeapi.co/api/v2/location/{id or name}/
	if data, ok := (*conf.Redis).Get(url); ok {
		var location Location

		if err := json.Unmarshal(data, &location); err != nil {
			return Location{}, err
		}

		return location, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Location{}, err
	}

	defer resp.Body.Close()

	var location Location

	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return Location{}, err
	}

	encoded, err := json.Marshal(location)

	if err == nil {
		(*conf.Redis).Set(url, encoded, 1*time.Hour)
	}

	return location, nil
}

func ApiGetLocationArea(url string, conf *config) (Location, error) {
	// GET https://pokeapi.co/api/v2/location-area/{id or name}/
	if data, ok := (*conf.Redis).Get(url); ok {
		var location Location

		if err := json.Unmarshal(data, &location); err != nil {
			return Location{}, err
		}

		return location, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Location{}, err
	}

	defer resp.Body.Close()

	var location Location

	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return Location{}, err
	}

	encoded, err := json.Marshal(location)

	if err == nil {
		(*conf.Redis).Set(url, encoded, 1*time.Hour)
	}

	return location, nil
}

func ApiGetAreaDetails(url string, conf *config) (LocationArea, error) {
	// GET https://pokeapi.co/api/v2/location-area/{id or name}/
	if data, ok := (*conf.Redis).Get(url); ok {
		var locationArea LocationArea

		if err := json.Unmarshal(data, &locationArea); err != nil {
			return LocationArea{}, err
		}

		return locationArea, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}

	defer resp.Body.Close()

	var locationArea LocationArea

	if err := json.NewDecoder(resp.Body).Decode(&locationArea); err != nil {
		return LocationArea{}, err
	}

	encoded, err := json.Marshal(locationArea)

	if err == nil {
		(*conf.Redis).Set(url, encoded, 1*time.Hour)
	}

	return locationArea, nil
}

func ApiGetPokemon(url string, conf *config) (Pokemon, error) {
	// GET https://pokeapi.co/api/v2/pokemon/{id or name}/
	if data, ok := (*conf.Redis).Get(url); ok {
		var pokemon Pokemon

		if err := json.Unmarshal(data, &pokemon); err != nil {
			return Pokemon{}, err
		}

		return pokemon, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}

	defer resp.Body.Close()

	var pokemon Pokemon

	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return Pokemon{}, err
	}

	encoded, err := json.Marshal(pokemon)

	if err == nil {
		(*conf.Redis).Set(url, encoded, 1*time.Hour)
	}

	return pokemon, nil
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)

	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range conf.Commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMapForward(conf *config) error {
	if conf.Map.Next == "" {
		fmt.Println("You're on the last page")
		return nil
	}

	location, err := ApiGetLocationArea(conf.Map.Next, conf)
	if err != nil {
		return err
	}

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}

	conf.Map.Next = location.Next
	conf.Map.Previous = location.Previous

	return nil
}

func commandMapBack(conf *config) error {
	if conf.Map.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	location, err := ApiGetLocationArea(conf.Map.Previous, conf)
	if err != nil {
		return err
	}

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}

	conf.Map.Next = location.Next
	conf.Map.Previous = location.Previous

	return nil
}

func commandExploreArea(conf *config) error {
	if len(conf.Args) < 1 {
		fmt.Println("Usage: explore <area>")

		return nil
	}

	area := conf.Args[0]

	areaDetails, err := ApiGetAreaDetails("https://pokeapi.co/api/v2/location-area/"+area, conf)

	if err != nil {
		return err
	}

	fmt.Println("Exploring " + areaDetails.Name)

	if len(areaDetails.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")

		for _, encounter := range areaDetails.PokemonEncounters {
			fmt.Println(" - " + encounter.Pokemon.Name)
		}
	} else {
		fmt.Println("No Pokemon found in this area")
	}

	return nil
}

func commandCatchPokemon(c *config) error {
	if len(c.Args) < 1 {
		fmt.Println("Usage: catch <pokemon>")

		return nil
	}

	pokemon := c.Args[0]

	fmt.Println("Throwing a Pokeball at " + pokemon + "...")

	pokemonDetails, err := ApiGetPokemon("https://pokeapi.co/api/v2/pokemon/"+pokemon, c)

	if err != nil {
		return err
	}

	baseExp := pokemonDetails.BaseExperience

	fmt.Println("Base experience: " + fmt.Sprint(baseExp))

	if rand.Intn(baseExp) < baseExp/4 {
		fmt.Println(pokemon + " was caught!")

		c.Pokemon[pokemon] = pokemonDetails
	}

	return nil
}

func commandInspectPokemon(c *config) error {
	if len(c.Args) < 1 {
		fmt.Println("Usage: inspect <pokemon>")

		return nil
	}

	pokemon := c.Args[0]

	if details, ok := c.Pokemon[pokemon]; ok {
		fmt.Println("Name: " + details.Name)

		fmt.Println("Height: " + fmt.Sprint(details.Height))
		fmt.Println("Weight: " + fmt.Sprint(details.Weight))

		fmt.Println("Stats:")

		for _, stat := range details.Stats {
			fmt.Println(" -" + stat.Stat.Name + ": " + strconv.Itoa(stat.BaseStat))
		}

		fmt.Println("Types:")

		for _, type_ := range details.Types {
			fmt.Println(" - " + type_.Type.Name)
		}
	} else {
		fmt.Println(pokemon + " has not been caught")
	}

	return nil
}

func main() {
	var input string

	var conf config

	conf.Redis = new(Caches)

	conf.Commands = make(map[string]cliCommand)

	conf.Map.Next = "https://pokeapi.co/api/v2/location-area/"
	conf.Map.Previous = ""

	conf.Commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    func() error { return commandHelp(&conf) },
	}

	conf.Commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    func() error { return commandExit(&conf) },
	}

	conf.Commands["map"] = cliCommand{
		name:        "map",
		description: "This will be how we explore the Pokemon world",
		callback:    func() error { return commandMapForward(&conf) },
	}

	conf.Commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "This will be how we explore the Pokemon world backwards",
		callback:    func() error { return commandMapBack(&conf) },
	}

	conf.Commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explore the map of a single area",
		callback: func() error {
			return commandExploreArea(&conf)
		},
	}

	conf.Commands["catch"] = cliCommand{
		name:        "catch",
		description: "Catch a Pokemon",
		callback: func() error {
			return commandCatchPokemon(&conf)
		},
	}

	conf.Commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Inspect a Pokemon",
		callback: func() error {
			return commandInspectPokemon(&conf)
		},
	}

	conf.Commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback: func() error {
			fmt.Println("Your Pokedex:")

			for name := range conf.Pokemon {
				fmt.Println(" - " + name)
			}

			return nil
		},
	}

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))

	fmt.Println("Welcome to the Pokedex!")

	// Our collection of Pokemon
	conf.Pokemon = make(map[string]Pokemon)

	// REPL loop
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		input = scanner.Text()

		clean := cleanInput(input)

		fmt.Println("Your command was: " + clean[0])

		if command, ok := conf.Commands[clean[0]]; ok {
			conf.Args = clean[1:]

			if err := command.callback(); err != nil {
				fmt.Println("Error: " + err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
