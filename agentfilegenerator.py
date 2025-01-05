import json
import random
import numpy as np

# Constants
SUBTYPES = ["None", "Pirate", "Converter"]
RELATION_TYPES = {
    "1": 0.75,  # Ennemi
    "2": 1.0,   # Pas de lien direct
    "3": 1.25,  # Amis
    "4": 1.5    # Famille
}

def generate_agents(num_believers, num_sceptics, num_neutrals, relations, random_relations=False, distribution="uniform", params=None, personal_param_distribution="uniform", personal_param_params=None):
    agents = []
    total_agents = num_believers + num_sceptics + num_neutrals

    for i in range(total_agents):
        agent_id = f"Agent{i}"
        if i < num_believers:
            opinion = round(random.uniform(2/3, 1), 2)
            agent_type = "Believer"
        elif i < num_believers + num_sceptics:
            opinion = round(random.uniform(0, 1/3), 2)
            agent_type = "Sceptic"
        else:
            opinion = round(random.uniform(1/3, 2/3), 2)
            agent_type = "Neutral"

        if distribution == "uniform":
            charisme = {f"Agent{j}": round(random.uniform(0, 1), 2) for j in range(total_agents) if j != i}
        elif distribution == "normal":
            mean = params.get("mean", 0.5)
            std_dev = params.get("std_dev", 0.1)
            charisme = {f"Agent{j}": round(np.clip(np.random.normal(mean, std_dev), 0, 1), 2) for j in range(total_agents) if j != i}
        
        if random_relations:
            relation = {f"Agent{j}": round(random.choice([0.75, 1.0, 1.25, 1.5]), 2) for j in range(total_agents) if j != i}
        else:
            relation = {f"Agent{j}": relations[agent_type] for j in range(total_agents) if j != i}
        
        if personal_param_distribution == "uniform":
            personal_parameter = round(random.uniform(personal_param_params["min"], personal_param_params["max"]), 2)
        elif personal_param_distribution == "normal":
            mean = personal_param_params.get("mean", 3.0)
            std_dev = personal_param_params.get("std_dev", 0.1)
            personal_parameter = round(np.clip(np.random.normal(mean, std_dev), 0, 6), 2)
        
        sub_type = random.choice(SUBTYPES)
        
        agent = {
            "id": agent_id,
            "opinion": opinion,
            "charisme": charisme,
            "relation": relation,
            "personalParameter": personal_parameter,
            "subType": sub_type
        }
        agents.append(agent)
    return agents

def save_agents_to_file(agents, filename):
    with open(filename, 'w') as file:
        json.dump(agents, file, indent=4)

if __name__ == "__main__":
    num_believers = int(input("Enter the number of Believers: "))
    num_sceptics = int(input("Enter the number of Sceptics: "))
    num_neutrals = int(input("Enter the number of Neutrals: "))
    
    random_relations = input("Do you want random relations between agents? (yes/no): ").lower() == "yes"
    
    relations = {}
    if not random_relations:
        for agent_type in ["Believer", "Sceptic", "Neutral"]:
            print(f"Choose relation type for {agent_type}s with other agents:")
            print("1 - Ennemi")
            print("2 - Pas de lien direct")
            print("3 - Amis")
            print("4 - Famille")
            choice = input(f"Relation type (1-4) for {agent_type}s: ")
            relations[agent_type] = RELATION_TYPES.get(choice, 1.0)
    
    distribution = input("Choose the distribution for charisma (uniform/normal): ").lower()
    params = {}
    if distribution == "normal":
        params["mean"] = float(input("Enter the mean for the normal distribution: "))
        params["std_dev"] = float(input("Enter the standard deviation for the normal distribution: "))

    personal_param_distribution = input("Choose the distribution for personalParameter (uniform/normal): ").lower()
    personal_param_params = {}
    if personal_param_distribution == "uniform":
        personal_param_params["min"] = float(input("Enter the minimum value for the uniform distribution: "))
        personal_param_params["max"] = float(input("Enter the maximum value for the uniform distribution: "))
    elif personal_param_distribution == "normal":
        personal_param_params["mean"] = float(input("Enter the mean for the normal distribution: "))
        personal_param_params["std_dev"] = float(input("Enter the standard deviation for the normal distribution: "))

    agents = generate_agents(num_believers, num_sceptics, num_neutrals, relations, random_relations, distribution, params, personal_param_distribution, personal_param_params)
    save_agents_to_file(agents, "agents.json")
    print(f"Generated {len(agents)} agents and saved to agents.json")
