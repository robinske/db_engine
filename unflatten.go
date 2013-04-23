var cacheData = Dictionary {}
type Dictionary map[string]interface{}

func unflatten(jsonString string) {
	for k, v := range cacheData {
		if strings.Contains(k, ":") {
			fmt.Println("split that shit")
		} else {
			fmt.Println("insert this into the string")
		}
	}
}


"""
if key does not contain ':'
	key:value in the top level of json
else:	
	split it
	if split value can be converted to a number, then it's a list index
"""



key:value in the top level of json

"""

batters:batter:1:type:Chocolate
id:0001
topping:2:id:5005
topping:0:id:5001
batters:batter:0:type:Regular
batters:batter:0:id:1001
batters:batter:2:type:Blueberry
type:donut
topping:6:id:5004
topping:5:id:5003
ppu:0.55
batters:batter:1:id:1002

topping:1:type:Glazed
batters:batter:3:type:Devil's Food
batters:batter:2:id:1003
name:Cake
topping:3:type:Powdered Sugar
topping:1:id:5002
topping:5:type:Chocolate
topping:6:type:Maple
topping:4:type:Chocolate with Sprinkles
topping:4:id:5006
topping:0:type:None
topping:2:type:Sugar
batters:batter:3:id:1004
topping:3:id:5007

"""



{
	"id": "0001",
	"type": "donut",
	"name": "Cake",
	"ppu": 0.55,
	"batters":{"batter":[{ "id": "1001", "type": "Regular" },{ "id": "1002", "type": "Chocolate" },{ "id": "1003", "type": "Blueberry" },{ "id": "1004", "type": "Devil's Food" }]
		},
	"topping":
		[
			{ "id": "5001", "type": "None" },
			{ "id": "5002", "type": "Glazed" },
			{ "id": "5005", "type": "Sugar" },
			{ "id": "5007", "type": "Powdered Sugar" },
			{ "id": "5006", "type": "Chocolate with Sprinkles" },
			{ "id": "5003", "type": "Chocolate" },
			{ "id": "5004", "type": "Maple" }
		]
}