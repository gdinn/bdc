# Projeto BDC - Base de Dados Condominial

## 📋 Sobre o Projeto

### Objetivos
- **Primário:** Desenvolvimento e aprimoramento de habilidades em back-end utilizando Go
- **Secundário:** Prática com infraestrutura e pipelines de CI/CD

### Descrição
O BDC é uma API back-end desenvolvida para gerenciamento de informações condominiais. O sistema permite que moradores e administradores gerenciem dados relacionados aos apartamentos de um condomínio através de uma interface front-end.

## 🚀 Funcionalidades

O sistema possibilita aos usuários:
- Inserir informações sobre apartamentos [Desenvolvendo]
- Gerenciar dados condominiais [Backlog]
- Comunicação entre front-end e back-end via API REST [Backlog]

## 📄 Licença

Este projeto está licenciado sob a **MIT License**. Consulte o arquivo [LICENSE](https://github.com/gdinn/bdc/blob/main/LICENSE) para mais informações.
```
bdc
├─ 📁back-end
│  ├─ 📁api
│  │  └─ 📄routes.go
│  ├─ 📁cmd
│  │  └─ 📄main.go
│  ├─ 📁cognito
│  │  └─ 📄pre-token-generation.js # Lambda de geração de token do cognito
│  ├─ 📁configs
│  │  ├─ 📄.aws_config_example
│  │  └─ 📄.env_example.env
│  ├─ 📁internal
│  │  ├─ 📁database
│  │  │  ├─ 📄connection.go
│  │  │  ├─ 📄migration.sql
│  │  │  ├─ 📄migrations.go
│  │  │  ├─ 📄notes.md
│  │  │  └─ 📄setup.sql
│  │  ├─ 📁domain
│  │  │  ├─ 📄errors.go
│  │  │  └─ 📄user_context.go
│  │  ├─ 📁handlers
│  │  │  ├─ 📄user.requests.sh
│  │  │  └─ 📄user_handler.go
│  │  ├─ 📁middleware
│  │  │  └─ 📄auth_middleware.go
│  │  ├─ 📁models
│  │  │  ├─ 📄apartamento.go
│  │  │  ├─ 📄base.go
│  │  │  ├─ 📄bicycle.go
│  │  │  ├─ 📄pet.go
│  │  │  ├─ 📄user.go
│  │  │  ├─ 📄user_apartment.go
│  │  │  └─ 📄vehicle.go
│  │  ├─ 📁repositories
│  │  │  ├─ 📄cognito_repository.go
│  │  │  └─ 📄user_repository.go
│  │  ├─ 📁services
│  │  │  ├─ 📄cognito_service.go
│  │  │  └─ 📄user_service.go
│  │  └─ 📁utils
│  │     └─ 📄response.go
│  ├─ 📄go.mod
│  └─ 📄go.sum
├─ 📁docs
│  ├─ 📁architecture
│  │  └─ 📄arch_diagram.mermaid
│  └─ 📁uml
│     ├─ 📄README.md
│     ├─ 📄prompt.md
│     └─ 📄uml_diagram.mermaid
├─ 📁tools
│  └─ 📄claude_packager.py
├─ 📄.gitignore
├─ 📄LICENSE
└─ 📄README.md
```