package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samluiz/blog/common/slug"
)

type Router interface {
	Home(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
}

type router struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) Router {
	return &router{app}
}

type PostOutput struct {
	ID          int
	Title       string
	PublishedAt string
	Content 	 string
	Slug 		 string
}

var Posts = []PostOutput{
	{
		ID:          1,
		Title:       "How I've built my blog using Go + Htmx + TailwindCSS",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     "This is the content of the post",
		Slug:        slug.GenerateSlug("How I've built my blog using Go + Htmx + TailwindCSS", slug.GenerateSlugId()),
	},

	{
		ID:          2,
		Title:       "How to build an API using Go and Fiber",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     "This is the content of the post",
		Slug:        slug.GenerateSlug("How to build an API using Go and Fiber", slug.GenerateSlugId()),
	},

	{
		ID:          3,
		Title:       "Como eu penso o JPA?",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     `<p>Modelar e descrever relações entre entidades pode ser uma dificuldade para desenvolvedores iniciantes, mas no caso do
		JPA pode virar uma completa bagunça se não prestar a devida atenção. Quando aprendi esta tecnologia criei um pequeno
		guia pra me ajudar e me baseio nele até hoje quando estou com alguma dúvida. É apenas minha pequena contribuição para
		a comunidade Java!</p>
	<h2> <a> </a> OneToMany ou ManyToOne relationship </h2>
	<p>Regra: geralmente, o lado “many” é o dono da relação.</p>
	<p>Exemplo: Cliente e Pedido.</p>
	<p>Perguntas a se fazer:</p>
	<ul>
		<li>Um cliente pode ter quantos pedidos?</li>
		<li>Um pedido pode pertencer a quantos clientes?</li>
		<li>Eu preciso saber quais os pedidos do cliente?</li>
		<li>Eu preciso saber qual o cliente de cada pedido?</li>
		<li>Quem é o dono do relacionamento?</li>
	</ul>
	<p>Respostas:</p>
	<ul>
		<li>[ ] Um cliente pode ter vários pedidos. (OneToMany)</li>
		<li>[ ] Um pedido pode pertencer a apenas um cliente. (ManyToOne)</li>
		<li>[ ] Eu não preciso saber os pedidos de cada cliente, mas preciso saber o cliente de cada pedido. (Unidirecional)
		</li>
		<li>[ ] O pedido é o dono do relacionamento.</li>
	</ul>
	<p>Portanto, chegamos a conclusão de que um cliente pode ter vários pedidos, mas um pedido pertence a apenas um cliente.
	</p>
	<p>Na classe pedido, nós declaramos o atributo cliente com as anotações:</p>
	<ul>
		<li>@ManyToOne (pois vários pedidos (Many) podem ter apenas um cliente (ToOne)</li>
		<li>(name = “cliente_id”) para criar uma coluna no banco de dados com o id do cliente</li>
	</ul>
	<p>Classe Pedido:<br> </p>
	<pre><code><span>@ManyToOne</span> <span>@JoinColumn</span><span>(</span><span>name</span> <span>=</span> <span>“</span><span>cliente_id</span><span>”</span><span>)</span> <span>private</span> <span>Cliente</span> <span>cliente</span><span>;</span> </code></pre>
	<p>Na classe cliente, nós declaramos e INSTANCIAMOS uma coleção de pedidos com as anotações:</p>
	<ul>
		<li>@JsonIgnore (para a resposta em json da requisição HTTP ignorar esta lista e não criar um loop de pedidos com
			cliente e cliente com pedidos, gerando stack overflow, afinal nós não precisamos saber os itens). Uma alternativa ao
			JsonIgnore pode ser usar DTO’s.</li>
		<li>@OneToMany(mappedBy = “cliente”) pois um cliente (One) pode ter vários pedidos (ToMany), o parâmetro mappedBy
			serve para mapear o atributo da classe Pedido (deve ser exatamente igual ao atributo declarado na classe)</li>
	</ul>
	<p>Classe Cliente:<br> </p>
	<pre><code><span>@OneToMany</span><span>(</span><span>mappedBy</span> <span>=</span> <span>“</span><span>cliente</span><span>”</span><span>)</span> <span>private</span> <span>List</span><span>&lt;</span><span>Pedido</span><span>&gt;</span> <span>pedidos</span><span>;</span> </code></pre>
	<hr>
	<h2> <a> </a> ManyToMany relationship </h2>
	<p>Exemplo: Um produto pode pertencer a várias categorias, e uma categoria possui vários produtos.</p>
	<p>Perguntas a se fazer:</p>
	<ul>
		<li>Um produto pode pertencer a quantas categorias?</li>
		<li>Uma categoria pode possuir quantos produtos?</li>
		<li>Eu preciso saber quais produtos uma categoria possui?</li>
		<li>Eu preciso saber a quais categorias um produto pertence?</li>
		<li>Quem é o dono do relacionamento?</li>
	</ul>
	<p>Respostas:</p>
	<ul>
		<li>[ ] Um produto pode pertencer a varias categorias. (ManyToMany)</li>
		<li>[ ] Uma categoria pode possuir vários produtos. (ManyToMany)</li>
		<li>[ ] Eu preciso saber quais produtos uma categoria possui, e também saber a quais categorias um produto pertence.
			(Bidirecional)</li>
		<li>[ ] O produto é o dono do relacionamento.</li>
	</ul>
	<p>Portanto, chegamos a conclusão de que um produto pode pertencer a várias categorias, e uma categoria pode possuir
		vários produtos. Nesse caso, devemos criar uma tabela de associação usando o @JoinTable. Como o Produto possui essa
		anotação, é implícito que o Produto é o dono da relação.</p>
	<p>Na classe produto, declaramos e INSTANCIAMOS uma coleção de categorias com as anotações:</p>
	<ul>
		<li>@ManyToMany (pois vários produtos (Many) podem ser de várias categorias (ToMany)</li>
		<li>@JoinTable (para criar uma nova tabela)</li>
		<li>@JoinColumn (dentro do @JoinTable para criar as colunas)</li>
	</ul>
	<p>Classe Produto:<br> </p>
	<pre><code><span>@ManyToMany</span> <span>@JoinTable</span><span>(</span><span>name</span> <span>=</span> <span>tb_product_category</span><span>,</span> <span>joinColumns</span> <span>=</span> <span>@JoinColumn</span><span>(</span><span>name</span> <span>=</span> <span>product_id</span><span>),</span> <span>inverseJoinColumns</span> <span>=</span> <span>@JoinColumn</span><span>(</span><span>name</span> <span>=</span> <span>category_id</span><span>))</span> <span>private</span> <span>Set</span><span>&lt;</span><span>Category</span><span>&gt;</span> <span>categories</span> <span>=</span> <span>new</span> <span>HashSet</span><span>&lt;&gt;();</span> </code></pre>
	<p>Na classe category, declaramos e INSTANCIAMOS uma coleção de produtos com as anotações:</p>
	<ul>
		<li>@ManyToMany(mappedBy = “categories”) para mapear a coleção de categorias da classe Produto e indicar que o dono do
			relacionamento é o Produto.</li>
	</ul>
	<p>Classe Categoria:<br> </p>
	<pre><code><span>@ManyToMany</span><span>(</span><span>mappedBy</span> <span>=</span> <span>categories</span><span>)</span> <span>private</span> <span>Set</span><span>&lt;</span><span>Product</span><span>&gt;</span> <span>products</span> <span>=</span> <span>new</span> <span>HashSet</span><span>&lt;&gt;();</span> </code></pre>
	<hr>
	<h2> <a> </a> OneToOne relationship </h2>
	<p>Exemplo: Um usuário pode ter um perfil associado, e um perfil pertence a apenas um usuário.</p>
	<p>Perguntas a se fazer:</p>
	<ul>
		<li>Um usuário pode ter quantos perfis?</li>
		<li>Um perfil pode pertencer a quantos usuários?</li>
		<li>Eu preciso saber o perfil de cada usuário?</li>
		<li>Eu preciso saber o usuário associado a cada perfil?</li>
		<li>Quem é o dono do relacionamento?</li>
	</ul>
	<p>Respostas:</p>
	<ul>
		<li>[ ] Um usuário pode ter apenas um perfil. (OneToOne)</li>
		<li>[ ] Um perfil pode pertencer a apenas um usuário. (OneToOne)</li>
		<li>[ ] Eu preciso saber o perfil de cada usuário, e também o usuário associado a cada perfil. (Bidirecional)</li>
		<li>[ ] O usuário é o dono do relacionamento.</li>
	</ul>
	<p>Portanto, concluímos que um usuário pode ter apenas um perfil, e um perfil pertence a apenas um usuário. Nesse caso,
		o usuário será o dono do relacionamento.</p>
	<p>Na classe <strong><code>Usuario</code></strong>, declaramos o atributo <strong><code>perfil</code></strong> com as
		anotações:<br> </p>
	<pre><code> <span>@OneToOne</span><span>(</span><span>mappedBy</span> <span>=</span> <span>usuario</span><span>)</span> <span>private</span> <span>Perfil</span> <span>perfil</span><span>;</span> </code></pre>
	<p>Na classe <strong><code>Perfil</code></strong>, declaramos o atributo <strong><code>usuario</code></strong> com as
		anotações:<br> </p>
	<pre><code> <span>@OneToOne</span> <span>@JoinColumn</span><span>(</span><span>name</span> <span>=</span> <span>usuario_id</span><span>)</span> <span>private</span> <span>Usuario</span> <span>usuario</span><span>;</span> </code></pre>
	<p>Observe que a anotação <strong><code>@OneToOne</code></strong> é utilizada em ambas as classes para indicar a relação
		One-To-One. A classe <strong><code>Perfil</code></strong> possui a anotação <strong><code>@JoinColumn</code></strong>
		para criar uma coluna no banco de dados com o <strong><code>id</code></strong> do usuário, que será usado para mapear
		o relacionamento. A classe <strong><code>Usuario</code></strong> utiliza a anotação
		<strong><code>mappedBy</code></strong> no atributo <strong><code>perfil</code></strong> para indicar que o
		relacionamento é mapeado pelo atributo <strong><code>usuario</code></strong> na classe
		<strong><code>Perfil</code></strong>.</p>
	<p>Dessa forma, temos uma relação One-To-One entre <strong><code>Usuario</code></strong> e
		<strong><code>Perfil</code></strong>, em que um usuário pode ter apenas um perfil e um perfil pertence a apenas um
		usuário.</p>`,
		Slug:	slug.GenerateSlug("My math roadmap: from zero to hero", slug.GenerateSlugId()),
	},
}

func (r *router) Home(c *fiber.Ctx) error {

	return c.Render("pages/home", fiber.Map{
		"Posts": Posts,
		"PageTitle": "home",
	})
}

func (r *router) Post(c *fiber.Ctx) error {
	var post PostOutput

	for _, p := range Posts {
		if p.Slug == c.Params("slug") {
			post = p
			break
		}
	}

	return c.Render("pages/post", fiber.Map{
		"Post": post,
		"PageTitle": post.Slug,
	})
}