package main

import "sort"

func hasSkill(c Character, s string) bool {
	for _, k := range c.Skills {
		if k == s {
			return true
		}
	}
	return false
}

func learnSkill(c *Character, s string) bool {
	if hasSkill(*c, s) {
		return false
	}
	c.Skills = append(c.Skills, s)
	sort.Strings(c.Skills)
	return true
}
