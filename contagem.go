package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	s "strings";
)

type Palavra struct{
	word string
	ocorrencias []int
}

// Funcao que le o conteudo do arquivo e retorna um slice the string com todas as linhas do arquivo
func lerArquivo(caminhoDoArquivo string) ([]string, error) {
	arquivo, err := os.Open(caminhoDoArquivo)
	// Caso tenha encontrado algum erro ao tentar abrir o arquivo retorne o erro encontrado
	if err != nil {
		return nil, err
	}
	// Garante que o arquivo sera fechado apos o uso
	defer arquivo.Close()

	// Cria um scanner que le cada linha do arquivo
	var linhas []string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		//texto += scanner.Text()
		linhas = append(linhas, scanner.Text())
	}

	// Retorna as linhas lidas e um erro se ocorrer algum erro no scanner
	//fmt.Println("a: %s",scanner);
	return linhas, scanner.Err()
}


func imprimirLista(lista []string){
	for indice, linha := range lista {
		fmt.Println(indice, "." +linha+".")
	}
}

func separarPalavras(linha string)([]string){
	var palavras []string;
	palavras = s.Split(linha, " ");
	return tratarPalavras(palavras);
}

func tratarPalavras(lpalavras []string) ([] string){
	var palavrasTratadas []string;
	for _, linha := range lpalavras {
		palavra := linha
		if s.Contains(palavra,"(") {
			palavra = s.ToLower(s.Replace(palavra, "(","",-1)) 
		} 
		if s.Contains(palavra,")") {
			palavra = s.ToLower(s.Replace(palavra, ")","",-1))
		}
		if s.Contains(palavra,",") {
			palavra = s.ToLower(s.Replace(palavra, ",","",-1))
		}
		if s.Contains(palavra,".") {
			var divisoes []string = s.SplitAfter(palavra, ".")
			palavra = s.ToLower(s.Replace(divisoes[0], ".","",-1))
		}
		if s.Contains(palavra,";") {
			palavra = s.ToLower(s.Replace(palavra, ";","",-1))
		}
		if s.Contains(palavra,":") {
			palavra = s.ToLower(s.Replace(palavra, ":","",-1))
		}
		if(palavra != ""){
			palavrasTratadas = append(palavrasTratadas, s.ToLower(palavra))
		}
	}
	return palavrasTratadas;
}
func pesquisarPalavra(lista []Palavra, palavra string)(int){
	for i := 0; i < len(lista); i++{
		if(palavra == lista[i].word){
			return i
		}
	}
	return -1
}
func contarOcorrencias(lista []string, palavra string)([]int){
	var ocorrencias []int

	for i := 0; i < len(lista); i++{
		if(palavra == lista[i]){
			ocorrencias = append(ocorrencias, i)
		}
	}
	return ocorrencias
}

func escreverArquivoInvertido(listaPalavras []Palavra, documento string) error{
	arquivo, err := os.Create(documento)
	if err != nil {
		return err
	}
	defer arquivo.Close()
	escritor := bufio.NewWriter(arquivo)
	for i, _ := range listaPalavras {
		fmt.Fprint(escritor, listaPalavras[i].word)
		fmt.Fprint(escritor, " ")
		for j, _ := range listaPalavras[i].ocorrencias{
			fmt.Fprint(escritor, listaPalavras[i].ocorrencias[j])
			if(j != len(listaPalavras[i].ocorrencias) -1){
				fmt.Fprint(escritor, ", ")
			}
		}
		fmt.Fprint(escritor, "\n")
	}
	return escritor.Flush()
}

func gerarArquivoInvertido(documentoEntrada string, documentoSaida string){
	var conteudo []string
	var palavras []string
	var listaPalavras []Palavra
	conteudo, err := lerArquivo(documentoEntrada);
	if err != nil {
		log.Fatalf("Erro:")
	}
	//Juntar todas as palavras do texto;
	for i := 0; i < len(conteudo); i++ {
		palavras = append(palavras, separarPalavras(conteudo[i])...)
	}
	conteudo = nil;
	//Gerar estruturas
	for i := 0; i< len(palavras); i++{
		if (pesquisarPalavra(listaPalavras, palavras[i]) == -1){
			listaPalavras = append(listaPalavras, Palavra{palavras[i],contarOcorrencias(palavras,palavras[i])})
		}
	}
	escreverArquivoInvertido(listaPalavras,documentoSaida)
}

func retornarRelevancia(documentoInvertido string, pesquisa string) float64{
	var conteudo []string
	var palavras []Palavra
	var qtTotalPalavras int = 0
	conteudo, err := lerArquivo(documentoInvertido)
	if err != nil {
		log.Fatalf("Erro:")
	}
	for i := 0; i<len(conteudo); i++{
		word, ocor := tratarLinhaInvertida(conteudo[i])
		qtTotalPalavras += len(ocor)
		palavras = append(palavras, Palavra{word, ocor})
	}
	//Agora já temos a estrutura montada só pra fazer a pesquisa e o cálculo
	return calcularRelevancia(palavras, pesquisa, qtTotalPalavras)
}

func calcularRelevancia(palavras []Palavra, pesquisa string, qtTotalPalavras int) float64{
	i := pesquisarPalavra(palavras, pesquisa)
	if i == -1{
		return 0
	}
	return (float64(len(palavras[i].ocorrencias))/float64(qtTotalPalavras))
}

func tratarLinhaInvertida(linha string)(string, []int){
	var componentesLinha []string
	var palavra string
	var ocorrencias []int
	linha = s.Replace(linha, ",", "",-1)
	componentesLinha = s.SplitAfterN(linha, " ",-1);
	for i := 0; i<len(componentesLinha); i++{
		componentesLinha[i] = s.Replace(componentesLinha[i], " ", "", -1)
	}
	palavra = componentesLinha[0]
	for i := 1; i<len(componentesLinha); i++{
		pal, err := strconv.Atoi(componentesLinha[i])
		if err != nil{
			fmt.Println("Erro", err)
		}
		ocorrencias = append(ocorrencias, pal)
	}
	return palavra, ocorrencias
}
func main() {
	//Função de Gerar um arquivo invertido a partir de um documento específico
	gerarArquivoInvertido("documento1.txt", "arqInv1.txt")
	//Função de recuperar um arquivo invertido a partir de um documento específico e retornar a relevância de uma querry
	fmt.Println(retornarRelevancia("arqInv1.txt", "para"))
}