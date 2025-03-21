SET search_path TO test_schema;

CREATE FUNCTION test_schema.notify_new_packet() RETURNS TRIGGER
    LANGUAGE plpgsql AS $$
BEGIN
    PERFORM pg_notify('NEW_PACKET', row_to_json(NEW)::text);
RETURN NEW;
END;
$$;

CREATE TRIGGER packet_insert_trigger
    AFTER INSERT ON test_schema.test_network_flow
    FOR EACH ROW
    EXECUTE FUNCTION test_schema.notify_new_packet();
