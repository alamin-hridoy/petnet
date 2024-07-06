package drone

kind:  "pipeline"
type?: string
name:  #Identifier

trigger: #Trigger

workspace: close({
	path: string
})

image_pull_secrets: [...#Identifier]

steps: [...#Step]
#Step: {
	name: #Identifier
	depends_on: [...#Identifier]
	image: #Image
	settings: {
		[NAME=_]: string
	}
	volumes: [...{
		name: #Identifier
		path: string
	}]
	environment: {
		[NAME=#Identifier]: string | {
			from_secret: #Identifier
		}
	}
	commands: [...string]
	when: #Trigger
}

#Identifier: =~"[a-zA-Z_][a-zA-Z_-]*"
#Image:      string & !~" "
#Trigger: {
	event?: #Event
	status?: [...#Status]
	cron?:   #Filter
	branch?: #Filter
}
#Filter: [...#Identifier] | {include: [...#Identifier]} | {exclude: [...#Identifier]}
#Event:  "push" | "pull_request" | "tag" | "promote" | "rollback"
#Status: "success" | "failure"

volumes: [...#Volume]
#Volume: {
	name: #Identifier
	host?: {
		path: string
	}
	temp?: {}
}
