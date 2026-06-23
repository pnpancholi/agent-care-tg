package bot

var idx = 0

func GetSetBackMessage() string {
	idx = (idx + 1) % len(setbackMessages)
	return setbackMessages[idx]
}

var setbackMessages = []string{
	"It's okay to stumble. What truly matters is how you pick yourself up and learn from it. You've got this.",
	"Don't let a setback define your entire journey. View it as a detour, not a dead end. Keep moving forward.",
	"Progress isn't linear. This moment is a chance to pause, reflect, and come back even stronger. Your resilience is your superpower.",
	"Every great story has moments of challenge. This setback is just a chapter, not the whole book. Write the next one with renewed determination.",
	"A bump in the road doesn't mean the road is over. Take a breath, re-evaluate, and remember why you started. You're capable of more than you think.",
	"This isn't a failure, it's feedback. Use this experience to gain valuable insights and adjust your approach. Growth often comes from these moments.",
	"Be kind to yourself. Setbacks are part of being human. Acknowledge the challenge, then refocus on your next positive step.",
	"Your dedication isn't measured by perfection, but by your persistence. This challenge is a test of that persistence, and you'll overcome it.",
	"Sometimes, the biggest leaps forward follow a moment of retreat. Use this time to gather your strength and strategy. You're still on track.",
	"One step back, two steps forward. This setback doesn't diminish your efforts; it provides an opportunity to refine them. You're stronger than you think.",
}
