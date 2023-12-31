<h1>API CRUD</h1>
<p><strong>Passo 1:</strong> Baixe o projeto e no terminal, execute o comando:</p>
<br>
<h2>Windows</h2>
<p><code>iwr https://encore.dev/install.ps1 | iex</code></p>
<h2>Linux</h2>
<p><code>curl -L https://encore.dev/install.sh | bash</code></p>
<br>
<p><strong>Passo 2:</strong> Instale o Docker, ou caso já possua, inicie-o.</p>
<br>
<p><strong>Passo 3:</strong> No terminal, acesse a pasta do projeto com o comando <code>cd api-zpe-crud</code>.</p>
<br>
<p><strong>Passo 4:</strong> Execute o comando <code>encore run</code>. Esse comando criará um contêiner no Docker e configurará o ambiente de desenvolvimento local.</p>
<br>
<p><strong>Para executar as URLs disponibilizadas, siga as instruções detalhadas a seguir:</strong></p>
<a>Create: <code>https://staging-api-zpe-crud-gemi.encr.app/create/users</code></a>
<p>Estrutura:</p>
<pre>
{
    "ID": 0,
    "Name": "",
    "Email": "",
    "Role": ""
}
</pre>
<br>
<a>Read: <code>https://staging-api-zpe-crud-gemi.encr.app/read/users/:id</code></a>
<p>Substitua o <code>:id</code> na URL pelo ID que deseja buscar.</p>
<br>
<a>Update</a>
<p>Estrutura:</p>
<pre>
{
    "ID": 0,
    "Name": "",
    "Email": "",
    "Role": "",
    "NewRole": ""
}
</pre>
<br>
<a>Delete</a>
<p>Estrutura:</p>
<pre>
{
    "ID": 0,
    "Name": ""
}
</pre>
<br>
<a>List All Users</a>
<p>Apenas execute o endpoint e exibirá todos os detalhes dos usuários cadastrados.</p>
