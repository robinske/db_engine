import csv
import json

f = open('users.csv', 'r')
reader = csv.DictReader(f, delimiter="|", fieldnames = ( "id","age","gender","occupation","zipcode" ) )
out = json.dumps( [ row for row in reader ] )

fo = open('users.json', 'w')
fo.write(out)