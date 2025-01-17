CREATE TABLE IF NOT EXISTS dados_cota (
    id SERIAL PRIMARY KEY,
    data DATE,
    precoAbertura DECIMAL,
    precoMax DECIMAL,
    precoMin DECIMAL,
    precoUltimo DECIMAL,
    precoMedio DECIMAL,
    totalNegocios INT,
    qtdTitulos BIGINT,
    volTitulos BIGINT,
    nomeEmpresa VARCHAR(12),
    codNegocio VARCHAR(12)
)

CREATE TABLE IF NOT EXISTS arquivos_processados (
    arquivo VARCHAR PRIMARY KEY
)