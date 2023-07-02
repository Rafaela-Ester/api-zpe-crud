<h1>API CRUD</h1>
<p><strong>Passo 1:</strong> Baixe o projeto e no terminal, execute o comando:</p>
<br>
<h2>Windows</h2>
<p>iwr https://encore.dev/install.ps1 | iex</p>
<h2>Linux</h2>
<p>curl -L https://encore.dev/install.sh | bash</p>
<br>
<p><strong>Passo 2:</strong> Instale o Docker, ou caso já possua, inicie-o.</p>
<br>
<p><strong>Passo 3:</strong> No terminal, acesse a pasta do projeto com o comando "cd api-zpe-crud".</p>
<br>
<p><strong>Passo 4:</strong> Execute o comando "encore run". Esse comando criará um contêiner no Docker e configurará o ambiente de desenvolvimento local.</p>
<br>
<p><strong>Para executar as URLs disponibilizadas, siga as instruções detalhadas a seguir:</strong></p>
<a>Create:  https://staging-api-zpe-crud-gemi.encr.app/create/users</a>
<p>Estrutura</p>
<p>{
    "ID": 0,
    "Name": "",
    "Email": "",
    "Role": ""
}</p>
<br>
<a>Read: https://staging-api-zpe-crud-gemi.encr.app/read/users/:id</a>

 
 
