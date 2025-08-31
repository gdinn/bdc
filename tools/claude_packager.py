import os
import shutil
from pathlib import Path

def achatar_diretorios(diretorio_origem, diretorio_destino, arquivos_ignorar=None, prefixos_ignorar=None):
    """
    Copia todos os arquivos de um diretório e suas subpastas para um diretório de destino,
    removendo a estrutura de pastas (achatando).
    
    Args:
        diretorio_origem (str): Caminho do diretório de origem
        diretorio_destino (str): Caminho do diretório de destino
        arquivos_ignorar (list): Lista de nomes de arquivos para ignorar
        prefixos_ignorar (list): Lista de prefixos de arquivos para ignorar
    
    Returns:
        dict: Estatísticas da operação (arquivos copiados, conflitos, etc.)
    """
    
    # Define listas padrão se não fornecidas
    if arquivos_ignorar is None:
        arquivos_ignorar = []
    if prefixos_ignorar is None:
        prefixos_ignorar = []
    
    origem = Path(diretorio_origem)
    destino = Path(diretorio_destino)
    
    # Verifica se o diretório de origem existe
    if not origem.exists():
        raise FileNotFoundError(f"Diretório de origem não encontrado: {origem}")
    
    if not origem.is_dir():
        raise NotADirectoryError(f"O caminho não é um diretório: {origem}")
    
    # Cria o diretório de destino se não existir
    destino.mkdir(parents=True, exist_ok=True)
    
    def deve_ignorar_arquivo(nome_arquivo):
        """Verifica se um arquivo deve ser ignorado baseado nas regras definidas"""
        # Verifica se o nome exato está na lista de arquivos a ignorar
        if nome_arquivo in arquivos_ignorar:
            return True
        
        # Verifica se o arquivo começa com algum dos prefixos a ignorar
        for prefixo in prefixos_ignorar:
            if nome_arquivo.startswith(prefixo):
                return True
        
        return False
    
    # Estatísticas
    arquivos_copiados = 0
    conflitos_resolvidos = 0
    arquivos_ignorados = 0
    erros = []
    
    # Percorre recursivamente todos os arquivos
    for arquivo in origem.rglob('*'):
        if arquivo.is_file():
            # Verifica se o arquivo deve ser ignorado
            if deve_ignorar_arquivo(arquivo.name):
                arquivos_ignorados += 1
                print(f"Ignorado: {arquivo.relative_to(origem)}")
                continue
            
            try:
                # Nome do arquivo de destino
                nome_arquivo = arquivo.name
                arquivo_destino = destino / nome_arquivo
                
                # Trata conflitos de nomes (se já existe um arquivo com o mesmo nome)
                contador = 1
                nome_base = arquivo.stem
                extensao = arquivo.suffix
                
                while arquivo_destino.exists():
                    novo_nome = f"{nome_base}_{contador}{extensao}"
                    arquivo_destino = destino / novo_nome
                    contador += 1
                    if contador == 2:  # Primeira vez que houve conflito
                        conflitos_resolvidos += 1
                
                # Copia o arquivo
                shutil.copy2(arquivo, arquivo_destino)
                arquivos_copiados += 1
                
                print(f"Copiado: {arquivo.relative_to(origem)} → {arquivo_destino.name}")
                
            except Exception as e:
                erro_msg = f"Erro ao copiar {arquivo}: {str(e)}"
                erros.append(erro_msg)
                print(f"ERRO: {erro_msg}")
    
    # Retorna estatísticas
    resultado = {
        'arquivos_copiados': arquivos_copiados,
        'arquivos_ignorados': arquivos_ignorados,
        'conflitos_resolvidos': conflitos_resolvidos,
        'erros': len(erros),
        'lista_erros': erros
    }
    
    return resultado

# Exemplo de uso
if __name__ == "__main__":
    # Exemplo de como usar a função com filtros
    arquivos_ignorar = [
        "user.requests.sh",
        "README.md",
        "requirements.txt",
        "go.mod",
        "go.sum",
        "notes.md"
    ]
    
    prefixos_ignorar = [
        ".",  # Ignorar .gitignore, .DS_Store, etc...
        "_",  # Ignorar arquivos que começam com underscore
        "temp"  # Ignorar arquivos temporários
    ]
    
    try:
        stats = achatar_diretorios(
            diretorio_origem="../back-end",
            diretorio_destino="../claude",
            arquivos_ignorar=arquivos_ignorar,
            prefixos_ignorar=prefixos_ignorar
        )
        
        print("\n" + "="*50)
        print("RESUMO DA OPERAÇÃO:")
        print(f"Arquivos copiados: {stats['arquivos_copiados']}")
        print(f"Arquivos ignorados: {stats['arquivos_ignorados']}")
        print(f"Conflitos de nome resolvidos: {stats['conflitos_resolvidos']}")
        print(f"Erros encontrados: {stats['erros']}")
        
        if stats['erros'] > 0:
            print("\nErros detalhados:")
            for erro in stats['lista_erros']:
                print(f"  - {erro}")
                
    except Exception as e:
        print(f"Erro na execução: {e}")