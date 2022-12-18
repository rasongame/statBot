import json, sqlite3

con = sqlite3.connect("bot.db")
cur = con.cursor()
file = open("flood_bak.log")
for line in file.readlines():
    try:
        r_json = json.loads(line)
        chat_id = r_json["chat"]["id"]
        message_id = r_json["message_id"]
        user_id = r_json["from"]["id"]
        date = r_json['date']
        first_name = r_json['from']['first_name']
        lastname = ""
        username = ""
        if "username" in r_json['from']:
            username = r_json['from']['username']
        if "last_name" in r_json['from']:
            lastname = r_json['from']['last_name']
        text = None
        if "text" in r_json:
            text = r_json['text']
            cur.execute("INSERT INTO chat_messages VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
                        (chat_id, message_id, user_id, text, date, first_name, lastname, username))
            print("commit: " + str(message_id))
    except json.JSONDecodeError:
        pass

con.commit()
