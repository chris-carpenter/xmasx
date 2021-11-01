package xmasx

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"math/rand"
	"os"
	"strings"
)

type Person struct {
	Name      string
	Family    string
	GiftsTo   int    `default:"-1"`
	GiftsFrom int    `default:"-1"`
}

func (p Person) String() string {
	return fmt.Sprintf("%s of %s",p.Name,p.Family)
}

type GiftExList struct {
	People []Person
}

func (l GiftExList) String() string {
	output := ""
	for _, p := range l.People {
		output += fmt.Sprintf("%s gifts to %s\n",p,l.People[p.GiftsTo])
	}
	return output
}

func (l GiftExList) Families() map[string]int {
	families := make(map[string]int)
	for _, p := range l.People {
		if p.GiftsTo == -1 {
			families[p.Family] += 1
		}
	}
	return families
}

func (l GiftExList) getFamilyMembers(family string) []string {
	var members []string
	for _, p := range l.People {
		if p.Family == family {
			members = append(members, p.Name)
		}
	}
	return members
}

func (l GiftExList) getGiftSender(family string) int {
	senderList := funk.Filter(l.People, func(val Person) bool { return val.Family == family && val.GiftsTo == -1 }).([]Person)
	return funk.IndexOf(l.People, senderList[0])
}

func (l GiftExList) getRandomGiftRecipient(exclude string) int {
	recipientList := funk.Filter(l.People, func(val Person) bool { return val.Family != exclude && val.GiftsFrom == -1 }).([]Person)
	index := rand.Intn(len(recipientList))
	return funk.IndexOf(l.People, recipientList[index]) //This may need a func to better match Persons
}

func unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func readCsvFile(filePath string) GiftExList {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal().Err(err).Str("filepath", filePath).Msg("Unable to read input file")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal().Err(err).Str("filepath", filePath).Msg("Unable to parse file as CSV")
	}

	var people GiftExList
	for _, row := range records {
		p := &Person{}
		if err := defaults.Set(p); err != nil {
			log.Error().Err(err).Str("row",strings.Join(row," ")).Msg("Failed to set defaults")
		}
		p.Name = row[0]
		p.Family = row[1]
		people.People = append(people.People, *p)
	}

	return people
}

func Run(file string) {
	giftExList := readCsvFile(file)
	if len(giftExList.People)%2 == 1 {
		log.Fatal().Msg("Odd number of gift recipients")
	}
	families := RankByCount(giftExList.Families())
	for len(families) > 0 {
		largestFamily := families[0].Key
		senderIndex := giftExList.getGiftSender(largestFamily)
		recipientIndex := giftExList.getRandomGiftRecipient(largestFamily)
		giftExList.People[senderIndex].GiftsTo = recipientIndex
		giftExList.People[recipientIndex].GiftsFrom = senderIndex
		families = RankByCount(giftExList.Families())
	}
	val, err := json.MarshalIndent(giftExList, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("Unable to Marshal result")
	}
	fmt.Println(string(val))
	fmt.Println(giftExList)
}
