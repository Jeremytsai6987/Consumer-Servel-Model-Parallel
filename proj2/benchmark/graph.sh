#!/bin/bash
#
#SBATCH --mail-user=jeremyyawei@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/jeremyyawei/Parallel_Programming/project2/project-2-Jeremytsai6987/proj2/benchmark
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=60:00

module load python/3.9
module load golang/1.19
echo "Running Python script to generate speedup graphs..."
python3 plot.py
echo "Benchmark and plotting completed."
