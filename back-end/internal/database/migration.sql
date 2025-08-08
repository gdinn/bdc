-- =====================================================
-- Configuração de usuário para migrations - BDC
-- =====================================================

-- 1. Criar usuário específico para migrations
CREATE USER migration_user WITH
  LOGIN
  NOSUPERUSER
  NOCREATEDB
  NOCREATEROLE
  INHERIT
  NOREPLICATION
  CONNECTION LIMIT -1
  PASSWORD 'migration_password_changeit';

-- 1. Conectar no schema correto
SET search_path TO bdc;

-- 2. Transferir ownership de todas as tabelas existentes para migration_user
ALTER TABLE IF EXISTS users OWNER TO migration_user;
ALTER TABLE IF EXISTS apartments OWNER TO migration_user;
ALTER TABLE IF EXISTS user_apartments OWNER TO migration_user;
ALTER TABLE IF EXISTS vehicles OWNER TO migration_user;
ALTER TABLE IF EXISTS pets OWNER TO migration_user;
ALTER TABLE IF EXISTS bicycles OWNER TO migration_user;

-- 3. Transferir ownership das sequences (se existirem)
DO $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN 
        SELECT schemaname, sequencename 
        FROM pg_sequences 
        WHERE schemaname = 'bdc'
    LOOP
        EXECUTE format('ALTER SEQUENCE %I.%I OWNER TO migration_user', 
                      rec.schemaname, rec.sequencename);
    END LOOP;
END
$$;

-- 4. Transferir ownership das functions existentes
DO $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN 
        SELECT n.nspname as schema_name, p.proname as function_name,
               pg_catalog.pg_get_function_identity_arguments(p.oid) as args
        FROM pg_catalog.pg_proc p
        LEFT JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
        WHERE n.nspname = 'bdc'
    LOOP
        EXECUTE format('ALTER FUNCTION %I.%I(%s) OWNER TO migration_user', 
                      rec.schema_name, rec.function_name, rec.args);
    END LOOP;
END
$$;

-- 5. Transferir ownership das views
DO $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN 
        SELECT schemaname, viewname 
        FROM pg_views 
        WHERE schemaname = 'bdc'
    LOOP
        EXECUTE format('ALTER VIEW %I.%I OWNER TO migration_user', 
                      rec.schemaname, rec.viewname);
    END LOOP;
END
$$;

-- 6. Garantir que objetos futuros sejam criados com migration_user como owner
ALTER DEFAULT PRIVILEGES IN SCHEMA bdc 
    GRANT ALL PRIVILEGES ON TABLES TO migration_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA bdc 
    GRANT ALL PRIVILEGES ON SEQUENCES TO migration_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA bdc 
    GRANT ALL PRIVILEGES ON FUNCTIONS TO migration_user;

-- Conceder permissões ao usuário da aplicação - verificar se precisa msm
GRANT ALL PRIVILEGES ON SCHEMA bdc TO migration_user;
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA bdc TO migration_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA bdc TO migration_user;

-- Definir search_path padrão para o usuário
ALTER USER migration_user SET search_path TO bdc, public;    

-- =====================================================
-- Verificar permissões (queries úteis para debug)
-- =====================================================

-- Ver permissões do schema
-- SELECT * FROM information_schema.schema_privileges WHERE grantee IN ('app_user', 'migration_user');

-- Ver permissões das tabelas
-- SELECT * FROM information_schema.table_privileges WHERE grantee IN ('app_user', 'migration_user');

