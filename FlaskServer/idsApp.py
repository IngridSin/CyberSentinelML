from FlaskServer.db import DB

if __name__ == "__main__":
    db = DB()

    try:
        db.test_query()
        # db.listen_for_new_packets()
        db.close()
    except KeyboardInterrupt:
        db.close()
