Monte um diagrama de classes para um sistema de cadastro de um condomínio. Considere os seguintes aspectos e relações:

## Usuários
São aqueles que acessam o sistema e que podem possuem 2..N qualificações, conforme relação a seguir:
### Qualificadores de usuários
- Morador ou Externo (obrigatório)
- Adulto ou Criança (obrigatório)
- Síndico
- Conselheiro (somente se for morador)
- Representante legal

### Permissões elevadas
Usuários que tiverem a qualificação de síndico e conselheiro tem permissões elevadas, conforme relação a seguir:
- Síndico: Escrita e leitura para todos os dados envolvendo os apartamentos
- Conselho: Leitura para todos os dados envolvendo os apartamentos

### Permissões comuns
Um usuário pode estar associado a 0..N apartamentos, tendo acesso de leitura e escrita a todos que estiver associado.

### Apartamentos
Os apartamentos possuem relações com os usuários e com os ativos ligados a unidade. Dessa forma, cada apartamento terá:
- 0 a N usuários
- 0 a 2 carros 
- 0 a 2 motos
- 0 a N pets
- 0 a N bicicletas