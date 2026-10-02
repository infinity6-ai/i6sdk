/**
 * @fileoverview 
 */

//usuario pega o sqlitez, faz o open e recebe uma Sqliteconn
// logo que recebe o sqliteconn precisa ter um finally com o close()
//todos os params precisam respeitas o esquema de ter opts (dentro de opts tem os parametros efetivamente)
class SqliteConn {
    //Fecha conexao
    async close() {

    }

    //retorna resultados
    // params são utilizaods quando tem um ? na query, sqlite tem que resolver os binds com uma apipropria do sqlite
    async query(query, params) {

    }

    //não retorna dada, throw exception no erro
    // params são utilizaods quando tem um ? na query, sqlite tem que resolver os binds com uma apipropria do sqlite
    async update(query, params) {

    }

    //retorna um resultado (como o query)
    // sempre procura uma tabela que chama "namedqueries" no sqlite, pega a query com o name que recebeu e executa o método query com essa query e params recebidos
    async namedQuery(name, params) {

    }
}

/**
 * Sqlitez client.
 */
class Sqlitez {
    //returns an sqlconn
    //usuario usa o downloader para pegar o arquivo e passa aqui o caminho do OPFS
    async open(sqlitepath) {

    }
}

const sqlitez = new Sqlitez()

export { sqlitez };
