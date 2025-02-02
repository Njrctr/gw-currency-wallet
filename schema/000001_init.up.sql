CREATE TABLE IF NOT EXISTS users (
    id serial not null PRIMARY KEY,
    username varchar(255) not null unique,
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE IF NOT EXISTS wallets (
    id serial not null PRIMARY KEY,
    usd decimal(10,2) not null default '0.00',
    eur decimal(10,2) not null default '0.00',
    rub decimal(10,2) not null default '0.00'
);

CREATE TABLE IF NOT EXISTS users_wallets (
    id serial not null PRIMARY KEY,
    user_id int references users(id) on delete cascade not null,
    wallet_id int references wallets(id) on delete cascade not null
);

CREATE FUNCTION wallets_balance_check() RETURNS trigger AS $balance_ckeck$
    BEGIN
        -- Проверить что баланс не ухлдит ниже нуля
        IF NEW.usd < 0 OR NEW.eur < 0 OR NEW.rub < 0 THEN
            RAISE EXCEPTION 'Недостаточно средст для выполнения операции';
        END IF;

        RETURN NEW;
    END;
$balance_ckeck$ LANGUAGE plpgsql;
CREATE TRIGGER wallets_balance_check BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE PROCEDURE wallets_balance_check();