import numpy as np
import matplotlib.pyplot as plt
import ast
import os
import imageio

def read_file(filename):
    '''
    Read the file and return the data.
    '''
    with open(filename, "r") as f:
        file_str = f.read()
        data = ast.literal_eval(file_str)
        n_iter = data[0][0]
        grid = data[1]
        grid = np.array(grid)
    return (n_iter, grid)


def generate_plot(filename):
    '''
    Plot and save the data.
    '''
    _, data = read_file(filename)
    plt.clf()
    plt.imshow(data)
    plt.xlabel('X')
    plt.ylabel('Y')
    plt.colorbar(orientation='vertical')
    plt.title('Conjugate Gradient Plot')

    save_filename = '.' + filename.split('.')[1] + '.png'

    plt.savefig(save_filename)
    
    plt.close()

def generate_plots(path):
    '''
    Generate all the plots.
    '''
    for filename in os.listdir(path):
        if filename.endswith('.txt'):
            generate_plot(path + '/' + filename)
    

def generate_movie(path):
    '''
    Generate a movie from the saved images.
    '''
    img_paths = []

    for filename in os.listdir(path):
        if filename.endswith('.png') and filename != 'x.png' and filename != 'b.png':
            # Sort by the number in the filename.
            img_paths.append(path + '/' + filename)
    
    img_paths.sort(key=lambda x: int(x.split(path + '/x_')[1].split('.png')[0]))
    images = [imageio.imread(img_path) for img_path in img_paths]
    imageio.mimsave(path + '/movie.gif', images)

def remove_text_files(path):
    '''
    Remove the text files.
    '''
    for filename in os.listdir(path):
        if filename.endswith('.txt'):
            os.remove(path + '/' + filename)

def remove_temp_png_files(path):
    '''
    Remove the temporary png files.
    '''
    for filename in os.listdir(path):
        if filename.endswith('.png') and filename != 'x.png' and filename != 'b.png':
            os.remove(path + '/' + filename)


if __name__ == '__main__':
    output_path = './output'
    generate_plots(output_path)
    generate_movie(output_path)
    remove_text_files(output_path)
    remove_temp_png_files(output_path)
    