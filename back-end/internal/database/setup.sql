-- =====================================================
-- Script de criação das tabelas do BDC (Base de Dados Condominial)
-- Banco de dados: PostgreSQL
-- =====================================================

-- Criar db
CREATE DATABASE bdc
    WITH 
    OWNER = igdebastiani
    ENCODING = 'UTF8'
    LC_COLLATE = 'pt_BR.UTF-8'
    LC_CTYPE = 'pt_BR.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    TEMPLATE = template0;
-- Comentário do banco
COMMENT ON DATABASE bdc IS 'Base de Dados Condominial - Sistema de gestão residencial';
-- Conecta no banco
\c bdc;

-- Criar schema se não existir
CREATE SCHEMA IF NOT EXISTS bdc;
SET search_path TO bdc;

DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
      CREATE USER app_user WITH
        LOGIN
        NOSUPERUSER
        NOCREATEDB
        NOCREATEROLE
        INHERIT
        NOREPLICATION
        CONNECTION LIMIT -1
        PASSWORD 'changeit';
   END IF;
END
$$;

GRANT USAGE ON SCHEMA bdc TO app_user;
GRANT CREATE ON SCHEMA bdc TO app_user;
GRANT ALL PRIVILEGES ON SCHEMA bdc TO app_user;

SELECT current_schema();

-- =====================================================
-- 1. TABELA DE USUÁRIOS
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Campos específicos
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    birth_date TIMESTAMPTZ,
    type TEXT NOT NULL DEFAULT 'EXTERNAL' CHECK (type IN ('RESIDENT', 'EXTERNAL')),
    age_group TEXT NOT NULL DEFAULT 'ADULT' CHECK (age_group IN ('ADULT', 'CHILD')),
    role TEXT NOT NULL DEFAULT 'COMMON' CHECK (role IN ('COMMON', 'MANAGER', 'ADVISOR')),
    
    -- Constraints
    CONSTRAINT uq_users_email UNIQUE (email)
);

-- Índices para users
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_type ON users(type);
CREATE INDEX idx_users_role ON users(role);

-- =====================================================
-- 2. TABELA DE APARTAMENTOS
-- =====================================================
CREATE TABLE IF NOT EXISTS apartments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Campos específicos
    number VARCHAR(10) NOT NULL,
    building VARCHAR(10) NOT NULL,
    legal_representative_id UUID,
    
    -- Constraints
    CONSTRAINT fk_users_legal_rep_apartments 
        FOREIGN KEY (legal_representative_id) 
        REFERENCES users(id) 
        ON DELETE SET NULL
);

-- Índices para apartments
CREATE INDEX idx_apartments_deleted_at ON apartments(deleted_at);
CREATE INDEX idx_apartments_number_building ON apartments(number, building);
CREATE INDEX idx_apartments_legal_representative_id ON apartments(legal_representative_id);

-- =====================================================
-- 3. TABELA DE RELACIONAMENTO USUÁRIOS-APARTAMENTOS
-- =====================================================
CREATE TABLE IF NOT EXISTS user_apartments (
    user_id UUID NOT NULL,
    apartment_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    -- Primary key composta
    PRIMARY KEY (user_id, apartment_id),
    
    -- Foreign keys
    CONSTRAINT fk_user_apartments_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    CONSTRAINT fk_user_apartments_apartment 
        FOREIGN KEY (apartment_id) 
        REFERENCES apartments(id) 
        ON DELETE CASCADE
);

-- Índices para user_apartments
CREATE INDEX idx_user_apartments_user_id ON user_apartments(user_id);
CREATE INDEX idx_user_apartments_apartment_id ON user_apartments(apartment_id);

-- =====================================================
-- 4. TABELA DE VEÍCULOS
-- =====================================================
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Campos específicos
    plate VARCHAR(10) NOT NULL,
    model VARCHAR(100) NOT NULL,
    color VARCHAR(50) NOT NULL,
    year BIGINT CHECK (year >= 1900 AND year <= 2030),
    parking_number VARCHAR(10),
    type VARCHAR(20) NOT NULL CHECK (type IN ('CAR', 'MOTORCYCLE')),
    apartment_id UUID NOT NULL,
    
    -- Constraints
    CONSTRAINT uq_vehicles_plate UNIQUE (plate),
    CONSTRAINT fk_apartments_vehicles 
        FOREIGN KEY (apartment_id) 
        REFERENCES apartments(id) 
        ON DELETE CASCADE
);

-- Índices para vehicles
CREATE INDEX idx_vehicles_deleted_at ON vehicles(deleted_at);
CREATE INDEX idx_vehicles_plate ON vehicles(plate);
CREATE INDEX idx_vehicles_apartment_id ON vehicles(apartment_id);
CREATE INDEX idx_vehicles_type ON vehicles(type);

-- =====================================================
-- 5. TABELA DE PETS
-- =====================================================
CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Campos específicos
    name VARCHAR(100) NOT NULL,
    species VARCHAR(20) NOT NULL CHECK (species IN ('DOG', 'CAT', 'BIRD', 'RABBIT', 'OTHER')),
    breed VARCHAR(100),
    size VARCHAR(20),
    apartment_id UUID NOT NULL,
    
    -- Constraints
    CONSTRAINT fk_apartments_pets 
        FOREIGN KEY (apartment_id) 
        REFERENCES apartments(id) 
        ON DELETE CASCADE
);

-- Índices para pets
CREATE INDEX idx_pets_deleted_at ON pets(deleted_at);
CREATE INDEX idx_pets_apartment_id ON pets(apartment_id);
CREATE INDEX idx_pets_species ON pets(species);

-- =====================================================
-- 6. TABELA DE BICICLETAS
-- =====================================================
CREATE TABLE IF NOT EXISTS bicycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Campos específicos
    color VARCHAR(50) NOT NULL,
    model VARCHAR(100),
    identification_tag VARCHAR(50),
    apartment_id UUID NOT NULL,
    
    -- Constraints
    CONSTRAINT fk_apartments_bicycles 
        FOREIGN KEY (apartment_id) 
        REFERENCES apartments(id) 
        ON DELETE CASCADE
);

-- Índices para bicycles
CREATE INDEX idx_bicycles_deleted_at ON bicycles(deleted_at);
CREATE INDEX idx_bicycles_apartment_id ON bicycles(apartment_id);

-- =====================================================
-- 7. TRIGGER PARA ATUALIZAR updated_at
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Aplicar trigger em todas as tabelas
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_apartments_updated_at 
    BEFORE UPDATE ON apartments 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vehicles_updated_at 
    BEFORE UPDATE ON vehicles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pets_updated_at 
    BEFORE UPDATE ON pets 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bicycles_updated_at 
    BEFORE UPDATE ON bicycles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 8. VIEWS ÚTEIS
-- =====================================================

-- View para usuários ativos
CREATE OR REPLACE VIEW active_users AS
SELECT * FROM users WHERE deleted_at IS NULL;

-- View para apartamentos com informações completas
CREATE OR REPLACE VIEW apartments_with_details AS
SELECT 
    a.*,
    u.name as legal_representative_name,
    u.email as legal_representative_email,
    COUNT(DISTINCT ua.user_id) as resident_count,
    COUNT(DISTINCT v.id) FILTER (WHERE v.deleted_at IS NULL) as vehicle_count,
    COUNT(DISTINCT p.id) FILTER (WHERE p.deleted_at IS NULL) as pet_count,
    COUNT(DISTINCT b.id) FILTER (WHERE b.deleted_at IS NULL) as bicycle_count
FROM apartments a
LEFT JOIN users u ON a.legal_representative_id = u.id
LEFT JOIN user_apartments ua ON a.id = ua.apartment_id
LEFT JOIN vehicles v ON a.id = v.apartment_id
LEFT JOIN pets p ON a.id = p.apartment_id
LEFT JOIN bicycles b ON a.id = b.apartment_id
WHERE a.deleted_at IS NULL
GROUP BY a.id, u.name, u.email;

-- =====================================================
-- 9. COMENTÁRIOS NAS TABELAS E COLUNAS
-- =====================================================
COMMENT ON TABLE users IS 'Tabela de usuários do condomínio';
COMMENT ON COLUMN users.type IS 'Tipo de usuário: RESIDENT (morador) ou EXTERNAL (externo)';
COMMENT ON COLUMN users.age_group IS 'Faixa etária: ADULT (adulto) ou CHILD (criança)';
COMMENT ON COLUMN users.role IS 'Perfil de acesso: COMMON (normal), ADVISOR (conselheiro) ou MANAGER (síndico)';

COMMENT ON TABLE apartments IS 'Tabela de apartamentos do condomínio';
COMMENT ON COLUMN apartments.legal_representative_id IS 'ID do representante legal do apartamento';

COMMENT ON TABLE user_apartments IS 'Tabela de relacionamento entre usuários e apartamentos';

COMMENT ON TABLE vehicles IS 'Tabela de veículos dos moradores';
COMMENT ON COLUMN vehicles.type IS 'Tipo de veículo: CAR (carro) ou MOTORCYCLE (moto)';

COMMENT ON TABLE pets IS 'Tabela de animais de estimação dos moradores';
COMMENT ON COLUMN pets.species IS 'Espécie do animal: DOG, CAT, BIRD, RABBIT, OTHER';

COMMENT ON TABLE bicycles IS 'Tabela de bicicletas dos moradores';

-- =====================================================
-- CONCEDER PERMISSÕES FINAIS
-- =====================================================
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA bdc TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA bdc TO app_user;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA bdc TO app_user;

-- Permissões para futuras tabelas/sequências
ALTER DEFAULT PRIVILEGES IN SCHEMA bdc GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA bdc GRANT USAGE, SELECT ON SEQUENCES TO app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA bdc GRANT EXECUTE ON FUNCTIONS TO app_user;

-- =====================================================
-- VERIFICAÇÃO FINAL
-- =====================================================
SELECT 'Schema bdc criado com sucesso!' as status;
SELECT schemaname, tablename FROM pg_tables WHERE schemaname = 'bdc';


-- =====================================================
-- Bonus: Cheat Sheet
-- =====================================================
-- \l - Listar todos os bancos de dados
-- \dt - Listar tabelas do schema atual
-- \d nome_tabela - Descrever estrutura de uma tabela
-- \d+ - Versão detalhada da descrição
-- \dn - Listar schemas
-- \du - Listar usuários/roles
-- \df - Listar funções
-- \dv - Listar views
-- Apaga tudo
    -- SET search_path TO bdc;
    -- DROP SCHEMA IF EXISTS bdc CASCADE;
    -- dropdb bdc
    -- dropuser migration_user
    -- \i setup.sql
    -- \i migration.sql
-- Conecta banco
    -- \c bdc;
    -- SET search_path TO bdc;
-- Apaga usuários
    -- TRUNCATE TABLE bdc.users CASCADE;
