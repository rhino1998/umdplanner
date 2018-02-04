package main

var schemaString = `
	schema {
		query: Query
	}

	type Query{
		class(code: String!): Class
		classes(geneds: [GenEd!], minCredits: Int, maxCredits: Int): [Class!]!
	}

	type Professor{
		name: String!
	}

	type Class {
		code: String!
		title: String!
		credits: Int!
		description: String
		prerequisite: String
		restriction: String
		geneds: [GenEd!]!
		sections(): [Section!]!
	}

	type Section{
		code: String!
		meetings: [Meeting]!
	}

	type Meeting{
		room: String!
		building: String!
		duration: Duration!
	}

	scalar Time

	type Duration{
		start: Time!
		end: Time!
	}

	enum GenEd {
		FSAW
		FSPW
		FSOC
		FSMA
		FSAR
		DSNS
		DSHU
		DSSP
		SCIS
		DVUP
		DVCC
	}
`
